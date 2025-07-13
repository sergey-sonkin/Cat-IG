export interface VideoPrompt {
  id: string;
  text: string;
  theme: string;
  createdAt: Date;
}

export interface GeneratedVideo {
  id: string;
  promptId: string;
  videoUrl: string;
  duration: number;
  createdAt: Date;
}

export interface InstagramAccount {
  id: string;
  username: string;
  accessToken: string;
  isMainAccount: boolean;
  isActive: boolean;
}

export interface PostPerformance {
  postId: string;
  accountId: string;
  likes: number;
  comments: number;
  shares: number;
  views: number;
  engagementRate: number;
  postedAt: Date;
}

export interface OptimalTimeSlot {
  hour: number;
  dayOfWeek: number;
  performanceScore: number;
}