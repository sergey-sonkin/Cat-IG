import OpenAI from "openai";
import type { VideoPrompt } from "./types";

export class PromptGenerator {
  private openai: OpenAI;
  
  constructor(apiKey: string) {
    this.openai = new OpenAI({ apiKey });
  }

  private readonly themes = [
    "existential dread",
    "corporate middle management",
    "gen z slang misuse",
    "cryptocurrency obsession",
    "wellness influencer parody",
    "linkedin motivational posts",
    "artisanal everything",
    "sustainable living anxiety",
    "dating app failures",
    "work from home chaos"
  ];

  private readonly situations = [
    "realizes they've been eating the same cardboard for 3 years",
    "discovers their humans are just large, hairless cats",
    "starts a podcast about the futility of chasing laser dots",
    "becomes a life coach for other cats",
    "opens a meditation retreat for anxious house pets",
    "launches a startup selling cardboard boxes as premium furniture",
    "writes passive-aggressive emails to their food dispenser",
    "practices mindfulness while knocking things off tables",
    "develops an elaborate conspiracy theory about vacuum cleaners",
    "starts a support group for cats with imposter syndrome"
  ];

  async generatePrompt(): Promise<VideoPrompt> {
    const theme = this.themes[Math.floor(Math.random() * this.themes.length)];
    const situation = this.situations[Math.floor(Math.random() * this.situations.length)];
    
    const systemPrompt = `You are a creative director for post-ironic cat content. Generate absurd, slightly meta video prompts that combine internet culture with cat behavior. Keep it weird but family-friendly.`;
    
    const userPrompt = `Create a short video prompt (1-2 sentences) about a cat dealing with "${theme}" where the cat ${situation}. Make it absurd and slightly self-aware.`;

    try {
      const response = await this.openai.chat.completions.create({
        model: "gpt-4o-mini",
        messages: [
          { role: "system", content: systemPrompt },
          { role: "user", content: userPrompt }
        ],
        max_tokens: 150,
        temperature: 0.9
      });

      const promptText = response.choices[0]?.message?.content || "A cat stares judgmentally at the camera while questioning the meaning of existence.";
      
      return {
        id: crypto.randomUUID(),
        text: promptText,
        theme,
        createdAt: new Date()
      };
    } catch (error) {
      console.error("Failed to generate prompt:", error);
      
      return {
        id: crypto.randomUUID(),
        text: `A cat ${situation} while contemplating ${theme}, occasionally making direct eye contact with the camera to break the fourth wall.`,
        theme,
        createdAt: new Date()
      };
    }
  }

  async generateBatch(count: number = 5): Promise<VideoPrompt[]> {
    const prompts = await Promise.allSettled(
      Array(count).fill(null).map(() => this.generatePrompt())
    );
    
    return prompts
      .filter((result): result is PromiseFulfilledResult<VideoPrompt> => result.status === 'fulfilled')
      .map(result => result.value);
  }
}