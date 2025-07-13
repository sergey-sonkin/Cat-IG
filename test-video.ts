import { VideoGenerator } from "./src/video-generator";
import { PromptGenerator } from "./src/prompt-generator";

async function testVideoGeneration() {
  console.log("🧪 Testing video generation...");
  
  const provider = (process.env.VIDEO_PROVIDER as any) || "veo2";
  console.log(`Using provider: ${provider}`);
  
  const videoGen = new VideoGenerator(
    process.env.GEMINI_API_KEY!,
    process.env.REPLICATE_API_KEY!,
    process.env.VERTEX_API_KEY!,
    process.env.GOOGLE_PROJECT_ID!,
    provider
  );
  
  const promptGen = new PromptGenerator(process.env.OPENAI_API_KEY!);
  
  try {
    console.log("Generating a test prompt...");
    const prompt = await promptGen.generatePrompt();
    console.log(`Generated prompt: "${prompt.text}"`);
    
    console.log("Generating video from prompt...");
    const video = await videoGen.generateVideo(prompt);
    
    console.log("✅ Video generation successful!");
    console.log(`Video ID: ${video.id}`);
    console.log(`Video URL: ${video.videoUrl}`);
    console.log(`Duration: ${video.duration} seconds`);
    
    return video;
  } catch (error) {
    console.error("❌ Video generation failed:", error);
    throw error;
  }
}

async function testCustomPrompt() {
  console.log("\n🎨 Testing custom prompt...");
  
  const videoGen = new VideoGenerator(
    process.env.GEMINI_API_KEY!,
    process.env.REPLICATE_API_KEY!,
    process.env.VERTEX_API_KEY!,
    process.env.GOOGLE_PROJECT_ID!,
    (process.env.VIDEO_PROVIDER as any) || "veo2"
  );
  
  const customPrompt = "A sophisticated orange tabby cat wearing tiny glasses, sitting at a miniature desk, looking directly at the camera with an expression that says 'I've seen things you wouldn't believe.' The cat occasionally adjusts its glasses and sighs dramatically.";
  
  try {
    console.log(`Testing custom prompt: "${customPrompt}"`);
    const video = await videoGen.testGeneration(customPrompt);
    
    console.log("✅ Custom prompt test successful!");
    console.log(`Video URL: ${video.videoUrl}`);
    
    return video;
  } catch (error) {
    console.error("❌ Custom prompt test failed:", error);
    throw error;
  }
}

if (import.meta.main) {
  Promise.all([
    testVideoGeneration(),
    testCustomPrompt()
  ]).then(() => {
    console.log("\n🎉 All tests completed!");
  }).catch((error) => {
    console.error("\n💥 Tests failed:", error);
    process.exit(1);
  });
}