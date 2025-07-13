# AI Cat Content Generator

Automated system for generating post-ironic cat videos using AI and posting them to Instagram accounts for A/B testing.

## Concept

1. LLM generates quirky cat video prompts
2. Google's video generation API creates videos from prompts
3. Content gets posted to small test accounts at optimal times
4. Performance tracking identifies best-performing content
5. Top content gets reposted to main "bread winner" account

## Architecture

- **Prompt Generator**: Creates absurd cat video prompts
- **Video Generator**: Google Imagen Video API integration
- **Instagram Poster**: Handles posting and scheduling
- **Performance Tracker**: Analytics and engagement metrics
- **A/B Testing**: Routes best content to main account

## Setup

```bash
bun install
bun run dev
```

## Current Status

- [x] Project initialization and structure
- [ ] Instagram API integration and policy research
- [ ] Google video generation API setup
- [ ] LLM prompt generation system
- [ ] Content scheduling system
- [ ] Performance tracking and analytics
- [ ] A/B testing logic
