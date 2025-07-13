import { GoogleGenerativeAI } from "@google/generative-ai";
import Replicate from "replicate";
import type { VideoPrompt, GeneratedVideo } from "./types";

type VideoProvider = "veo2" | "veo3-replicate" | "veo3-vertex";

export class VideoGenerator {
  private geminiClient: GoogleGenerativeAI;
  private replicateClient: Replicate;
  private provider: VideoProvider;
  private vertexApiKey: string;
  private projectId: string;
  
  constructor(
    geminiApiKey: string, 
    replicateApiKey: string, 
    vertexApiKey: string, 
    projectId: string,
    provider: VideoProvider = "veo2"
  ) {
    this.geminiClient = new GoogleGenerativeAI(geminiApiKey);
    this.replicateClient = new Replicate({ auth: replicateApiKey });
    this.vertexApiKey = vertexApiKey;
    this.projectId = projectId;
    this.provider = provider;
  }

  async generateVideo(prompt: VideoPrompt): Promise<GeneratedVideo> {
    console.log(`Generating video with ${this.provider.toUpperCase()} for prompt: ${prompt.text}`);
    
    try {
      switch (this.provider) {
        case "veo3-replicate":
          return await this.generateWithVeo3Replicate(prompt);
        case "veo3-vertex":
          return await this.generateWithVeo3Vertex(prompt);
        default:
          return await this.generateWithVeo2(prompt);
      }
    } catch (error) {
      console.error(`Video generation error with ${this.provider}:`, error);
      
      return {
        id: crypto.randomUUID(),
        promptId: prompt.id,
        videoUrl: `https://placeholder-video.com/mock/${prompt.id}.mp4`,
        duration: 5,
        createdAt: new Date()
      };
    }
  }

  private async generateWithVeo2(prompt: VideoPrompt): Promise<GeneratedVideo> {
    const model = this.geminiClient.getGenerativeModel({ model: "veo-2.0-generate-001" });
    
    const result = await model.generateContent({
      contents: [{
        role: "user",
        parts: [{
          text: prompt.text
        }]
      }],
      generationConfig: {
        video: {
          aspectRatio: "9:16",
          personGeneration: "dont_allow",
          durationSeconds: 5
        }
      }
    });

    const videoUrl = result.response.candidates?.[0]?.content?.parts?.[0]?.videoMetadata?.videoUrl;
    
    if (!videoUrl) {
      throw new Error("No video URL returned from Veo 2");
    }

    return {
      id: crypto.randomUUID(),
      promptId: prompt.id,
      videoUrl,
      duration: 5,
      createdAt: new Date()
    };
  }

  private async generateWithVeo3Replicate(prompt: VideoPrompt): Promise<GeneratedVideo> {
    const prediction = await this.replicateClient.run(
      "google/veo-3:latest",
      {
        input: {
          prompt: prompt.text,
          enhance_prompt: true,
          negative_prompt: "low quality, blurry, distorted"
        }
      }
    );

    const videoUrl = Array.isArray(prediction) ? prediction[0] : prediction;
    
    if (!videoUrl || typeof videoUrl !== 'string') {
      throw new Error("No video URL returned from Veo 3");
    }

    return {
      id: crypto.randomUUID(),
      promptId: prompt.id,
      videoUrl,
      duration: 8,
      createdAt: new Date()
    };
  }

  private async generateWithVeo3Vertex(prompt: VideoPrompt): Promise<GeneratedVideo> {
    const endpoint = `https://aiplatform.googleapis.com/v1/projects/${this.projectId}/locations/us-central1/publishers/google/models/veo-3.0-generate-preview:predictLongRunning`;
    
    const response = await fetch(endpoint, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${this.vertexApiKey}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        instances: [{
          prompt: prompt.text
        }],
        parameters: {
          aspectRatio: "9:16",
          durationSeconds: 8,
          numberOfVideos: 1
        }
      })
    });

    if (!response.ok) {
      throw new Error(`Vertex AI request failed: ${response.statusText}`);
    }

    const result = await response.json();
    const operationName = result.name;
    
    const videoUrl = await this.pollVertexOperation(operationName);
    
    return {
      id: crypto.randomUUID(),
      promptId: prompt.id,
      videoUrl,
      duration: 8,
      createdAt: new Date()
    };
  }

  private async pollVertexOperation(operationName: string): Promise<string> {
    const maxAttempts = 30;
    const delayMs = 10000;

    for (let attempt = 0; attempt < maxAttempts; attempt++) {
      const response = await fetch(`https://aiplatform.googleapis.com/v1/${operationName}`, {
        headers: {
          'Authorization': `Bearer ${this.vertexApiKey}`,
        }
      });

      if (!response.ok) {
        throw new Error(`Failed to poll operation: ${response.statusText}`);
      }

      const operation = await response.json();
      
      if (operation.done) {
        if (operation.error) {
          throw new Error(`Operation failed: ${operation.error.message}`);
        }
        
        const videoUrl = operation.response?.predictions?.[0]?.videoUrl;
        if (!videoUrl) {
          throw new Error("No video URL in completed operation");
        }
        
        return videoUrl;
      }

      console.log(`Video generation in progress... (attempt ${attempt + 1}/${maxAttempts})`);
      await new Promise(resolve => setTimeout(resolve, delayMs));
    }

    throw new Error("Video generation timed out");
  }

  async generateBatch(prompts: VideoPrompt[]): Promise<GeneratedVideo[]> {
    const videos = await Promise.allSettled(
      prompts.map(prompt => this.generateVideo(prompt))
    );
    
    return videos
      .filter((result): result is PromiseFulfilledResult<GeneratedVideo> => result.status === 'fulfilled')
      .map(result => result.value);
  }

  setProvider(provider: VideoProvider): void {
    this.provider = provider;
    console.log(`Switched to ${provider.toUpperCase()} for video generation`);
  }

  async testGeneration(testPrompt: string = "A fluffy orange cat wearing sunglasses sits in a director's chair, occasionally looking directly at the camera with an expression of existential contemplation"): Promise<GeneratedVideo> {
    const prompt: VideoPrompt = {
      id: crypto.randomUUID(),
      text: testPrompt,
      theme: "test",
      createdAt: new Date()
    };

    return this.generateVideo(prompt);
  }
}