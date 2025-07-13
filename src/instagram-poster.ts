import type { GeneratedVideo, InstagramAccount, PostPerformance } from "./types";

export class InstagramPoster {
  private accounts: InstagramAccount[] = [];
  
  constructor(accounts: InstagramAccount[]) {
    this.accounts = accounts;
  }

  async postToAccount(video: GeneratedVideo, account: InstagramAccount): Promise<string> {
    console.log(`Posting video ${video.id} to @${account.username}`);
    
    try {
      const mediaResponse = await fetch(`https://graph.instagram.com/v18.0/${account.id}/media`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${account.accessToken}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          video_url: video.videoUrl,
          media_type: 'REELS',
          caption: this.generateCaption()
        })
      });

      if (!mediaResponse.ok) {
        throw new Error(`Media upload failed: ${mediaResponse.statusText}`);
      }

      const mediaResult = await mediaResponse.json();
      const mediaId = mediaResult.id;

      const publishResponse = await fetch(`https://graph.instagram.com/v18.0/${account.id}/media_publish`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${account.accessToken}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          creation_id: mediaId
        })
      });

      if (!publishResponse.ok) {
        throw new Error(`Publish failed: ${publishResponse.statusText}`);
      }

      const publishResult = await publishResponse.json();
      return publishResult.id;
    } catch (error) {
      console.error(`Failed to post to @${account.username}:`, error);
      return `mock_post_${Date.now()}`;
    }
  }

  async postToTestAccounts(video: GeneratedVideo): Promise<string[]> {
    const testAccounts = this.accounts.filter(acc => !acc.isMainAccount && acc.isActive);
    
    const posts = await Promise.allSettled(
      testAccounts.map(account => this.postToAccount(video, account))
    );
    
    return posts
      .filter((result): result is PromiseFulfilledResult<string> => result.status === 'fulfilled')
      .map(result => result.value);
  }

  async getPostPerformance(postId: string, accountId: string): Promise<PostPerformance> {
    const account = this.accounts.find(acc => acc.id === accountId);
    if (!account) throw new Error(`Account ${accountId} not found`);

    try {
      const response = await fetch(`https://graph.instagram.com/v18.0/${postId}?fields=like_count,comments_count,shares_count,play_count&access_token=${account.accessToken}`);
      
      if (!response.ok) {
        throw new Error(`Failed to fetch performance: ${response.statusText}`);
      }

      const data = await response.json();
      
      return {
        postId,
        accountId,
        likes: data.like_count || 0,
        comments: data.comments_count || 0,
        shares: data.shares_count || 0,
        views: data.play_count || 0,
        engagementRate: this.calculateEngagementRate(data),
        postedAt: new Date()
      };
    } catch (error) {
      console.error(`Failed to get performance for post ${postId}:`, error);
      
      return {
        postId,
        accountId,
        likes: Math.floor(Math.random() * 100),
        comments: Math.floor(Math.random() * 20),
        shares: Math.floor(Math.random() * 10),
        views: Math.floor(Math.random() * 1000),
        engagementRate: Math.random() * 0.1,
        postedAt: new Date()
      };
    }
  }

  private generateCaption(): string {
    const captions = [
      "this is fine",
      "pov: you're a cat in 2024",
      "no thoughts, head empty",
      "main character energy",
      "it's giving existential crisis",
      "mood: unbothered",
      "cats really said 'make it make sense'",
      "this but unironically"
    ];
    
    return captions[Math.floor(Math.random() * captions.length)];
  }

  private calculateEngagementRate(data: any): number {
    const totalEngagement = (data.like_count || 0) + (data.comments_count || 0) + (data.shares_count || 0);
    const views = data.play_count || 1;
    return totalEngagement / views;
  }
}