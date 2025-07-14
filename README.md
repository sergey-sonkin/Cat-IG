# AI Cat Content Generator

Automated system for generating post-ironic cat videos using AI and posting them to Instagram accounts for A/B testing.

Built with Go for better performance, concurrency, and single binary deployment.

## Concept

1. LLM generates quirky cat video prompts
2. Google's video generation API creates videos from prompts
3. Content gets posted to small test accounts at optimal times
4. Performance tracking identifies best-performing content
5. Top content gets reposted to main "bread winner" account

## Architecture

- **Prompt Generator**: Creates absurd cat video prompts using OpenAI
- **Video Generator**: Supports Veo 2 (Gemini), Veo 3 (Replicate), and Veo 3 (Vertex AI)
- **Instagram Poster**: Handles posting and scheduling
- **Performance Tracker**: Analytics and engagement metrics with thread safety
- **A/B Testing**: Routes best content to main account

## Setup

```bash
# Install dependencies
make install

# Copy environment template (optional - can use shell env vars)
cp .env.example .env
# Edit .env with your API keys, OR set them in your shell

# Build
make build

# Run full pipeline (requires Instagram API keys)
go run main_full.go types.go prompt_generator.go video_generator.go instagram_poster.go performance_tracker.go

# Test video generation only (recommended first)
go run main.go types.go prompt_generator.go video_generator.go
```

## Testing Video Generation

The simplified test script (`main.go`) will:
1. Generate 3 absurd cat prompts using OpenAI
2. Create a video from the first prompt 
3. Test a custom business cat prompt
4. Display video URLs and generation times

**Required environment variables:**
- `OPENAI_API_KEY` - for prompt generation
- `GEMINI_API_KEY` - if using veo2 provider  
- `REPLICATE_API_KEY` - if using veo3-replicate provider
- `VIDEO_PROVIDER` - set to `veo2`, `veo3-replicate`, or `veo3-vertex`

**Quick test:**
```bash
export VIDEO_PROVIDER=veo3-replicate  # or veo2
go run main.go types.go prompt_generator.go video_generator.go
```

## Video Providers

- **veo2**: Gemini API (most accessible, cheaper, 5s videos)
- **veo3-replicate**: Replicate API (~$0.75/second, includes audio, 8s videos)
- **veo3-vertex**: Vertex AI (newest, requires allowlist access, 8s videos)

Set `VIDEO_PROVIDER` in your .env file.

## Current Status

- [x] Project initialization and Go structure
- [x] LLM prompt generation system
- [x] Google video generation API setup (all 3 providers)
- [x] Instagram API integration
- [x] Performance tracking and analytics
- [x] A/B testing logic
- [ ] Cron job scheduling for automated posting
- [ ] Database for persistent storage
