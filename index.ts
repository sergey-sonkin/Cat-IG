import { PromptGenerator } from "./src/prompt-generator";
import { VideoGenerator } from "./src/video-generator";
import { InstagramPoster } from "./src/instagram-poster";
import { PerformanceTracker } from "./src/performance-tracker";

console.log("🐱 AI Cat Content Generator starting...");

async function main() {
  const promptGen = new PromptGenerator(process.env.OPENAI_API_KEY!);
  const videoGen = new VideoGenerator(process.env.GOOGLE_API_KEY!);
  const tracker = new PerformanceTracker();
  
  const testAccounts = [
    {
      id: "test1",
      username: "cat_vibes_1",
      accessToken: process.env.INSTA_TOKEN_1!,
      isMainAccount: false,
      isActive: true
    },
    {
      id: "test2", 
      username: "cat_vibes_2",
      accessToken: process.env.INSTA_TOKEN_2!,
      isMainAccount: false,
      isActive: true
    },
    {
      id: "main",
      username: "main_cat_account",
      accessToken: process.env.INSTA_TOKEN_MAIN!,
      isMainAccount: true,
      isActive: true
    }
  ];
  
  const poster = new InstagramPoster(testAccounts);
  
  const prompts = await promptGen.generateBatch(3);
  console.log("Generated prompts:", prompts.map(p => p.text));
  
  const videos = await videoGen.generateBatch(prompts);
  console.log("Generated videos:", videos.length);
  
  for (const video of videos) {
    const postIds = await poster.postToTestAccounts(video);
    console.log(`Posted video ${video.id} to ${postIds.length} test accounts`);
  }
}

if (import.meta.main) {
  main().catch(console.error);
}