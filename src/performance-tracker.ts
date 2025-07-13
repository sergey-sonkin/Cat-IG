import type { PostPerformance, OptimalTimeSlot } from "./types";

export class PerformanceTracker {
  private performances: PostPerformance[] = [];
  
  addPerformance(performance: PostPerformance): void {
    this.performances.push(performance);
  }

  getBestPerformingPosts(limit: number = 10): PostPerformance[] {
    return this.performances
      .sort((a, b) => b.engagementRate - a.engagementRate)
      .slice(0, limit);
  }

  getOptimalPostingTimes(): OptimalTimeSlot[] {
    const timeSlots = new Map<string, { totalEngagement: number, count: number }>();
    
    this.performances.forEach(perf => {
      const hour = perf.postedAt.getHours();
      const dayOfWeek = perf.postedAt.getDay();
      const key = `${dayOfWeek}-${hour}`;
      
      if (!timeSlots.has(key)) {
        timeSlots.set(key, { totalEngagement: 0, count: 0 });
      }
      
      const slot = timeSlots.get(key)!;
      slot.totalEngagement += perf.engagementRate;
      slot.count += 1;
    });

    return Array.from(timeSlots.entries())
      .map(([key, data]) => {
        const [dayOfWeek, hour] = key.split('-').map(Number);
        return {
          hour,
          dayOfWeek,
          performanceScore: data.totalEngagement / data.count
        };
      })
      .sort((a, b) => b.performanceScore - a.performanceScore)
      .slice(0, 5);
  }

  shouldPromoteToMain(postPerformance: PostPerformance): boolean {
    const avgEngagement = this.getAverageEngagementRate();
    const threshold = avgEngagement * 1.5;
    
    return postPerformance.engagementRate > threshold && 
           postPerformance.views > 500;
  }

  getAverageEngagementRate(): number {
    if (this.performances.length === 0) return 0;
    
    const total = this.performances.reduce((sum, perf) => sum + perf.engagementRate, 0);
    return total / this.performances.length;
  }

  getAnalytics() {
    const totalPosts = this.performances.length;
    const totalViews = this.performances.reduce((sum, perf) => sum + perf.views, 0);
    const totalLikes = this.performances.reduce((sum, perf) => sum + perf.likes, 0);
    const avgEngagement = this.getAverageEngagementRate();
    
    return {
      totalPosts,
      totalViews,
      totalLikes,
      averageEngagementRate: avgEngagement,
      bestPerformingPost: this.getBestPerformingPosts(1)[0],
      optimalTimes: this.getOptimalPostingTimes()
    };
  }
}