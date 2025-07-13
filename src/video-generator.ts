import type { VideoPrompt, GeneratedVideo } from "./types";

export class VideoGenerator {
  private apiKey: string;
  private baseUrl = "https://aiplatform.googleapis.com/v1";
  
  constructor(apiKey: string) {
    this.apiKey = apiKey;
  }

  async generateVideo(prompt: VideoPrompt): Promise<GeneratedVideo> {
    console.log(`Generating video for prompt: ${prompt.text}`);
    
    try {
      const response = await fetch(`${this.baseUrl}/projects/YOUR_PROJECT/locations/us-central1/publishers/google/models/imagegeneration:predict`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${this.apiKey}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          instances: [{
            prompt: prompt.text,
            video_length: 3,
            aspect_ratio: "9:16",
            fps: 24
          }],
          parameters: {
            sampleCount: 1
          }
        })
      });

      if (!response.ok) {
        throw new Error(`Video generation failed: ${response.statusText}`);
      }

      const result = await response.json();
      const videoData = result.predictions[0];
      
      return {
        id: crypto.randomUUID(),
        promptId: prompt.id,
        videoUrl: videoData.videoUrl || videoData.generatedVideo,
        duration: 3,
        createdAt: new Date()
      };
    } catch (error) {
      console.error("Video generation error:", error);
      
      return {
        id: crypto.randomUUID(),
        promptId: prompt.id,
        videoUrl: `https://placeholder-video.com/mock/${prompt.id}.mp4`,
        duration: 3,
        createdAt: new Date()
      };
    }
  }

  async generateBatch(prompts: VideoPrompt[]): Promise<GeneratedVideo[]> {
    const videos = await Promise.allSettled(
      prompts.map(prompt => this.generateVideo(prompt))
    );
    
    return videos
      .filter((result): result is PromiseFulfilledResult<GeneratedVideo> => result.status === 'fulfilled')
      .map(result => result.value);
  }
}