//interface
export interface Images {
    small: string;
    grid: string;
    large: string;
    medium: string;
}

export type InfoboxItem = {
    key: string;
    value: string | Array<{ k?: string; v: string }>;
};

//单个条目结构
export interface BangumiSubject {
    id: number;
    name: string;
    summary: string;
    images: Images;
    infobox: InfoboxItem[];
    gender: string | null;
    stat: {
        comments: number;
        collects: number;
    }
}

//calendar
// 2. 评分统计对象
export interface BangumiRating {
  total: number;
  count: Record<string, number>; // 键是 "1"-"10"，值是人数
  score: number;
}

// 3. 收藏状态对象 (在看人数)
export interface BangumiCollection {
  doing: number;
}

// 4. 核心：每一个动画条目 (Item)
export interface BangumiSubjectCalendar {
  id: number;
  url: string;
  type: number;
  name: string;
  name_cn: string; // 有时为空字符串
  summary: string; // 有时为空
    air_date: string; // "YYYY-MM-DD"
    week_day: WeekdayInfo;
    
  
  // 注意：新番可能没有评分或排名，所以加 ? 号变成可选
    rating?: BangumiRating;

  rank?: number;
  collection?: BangumiCollection;
}

// 5. 星期几的信息
export interface WeekdayInfo {
  en: string; // "Mon"
  cn: string; // "星期一"
  ja: string; // "月耀日"
  id: number; // 1-7
}

// 6. 顶层：每一天的数据对象
export interface CalendarDay {
  weekday: WeekdayInfo;
  items: BangumiSubjectCalendar[];
}

// 7. 最终类型：API 返回的是一个数组
export type CalendarResponse = CalendarDay[];



/*
//////////////条目///////////////////
*/

export interface BangumiImages {
  small: string;
  grid: string;
  large: string;
  medium: string;
  common: string;
}

// 2. 标签对象
export interface BangumiTag {
  name: string;
  count: number;
  total_cont?: number; // 有些旧数据可能没有这个字段
}

// 4. 评分统计对象
export interface BangumiRating {
  rank: number;
  total: number;
  // 键是 "1" 到 "10"，值是打分人数
  count: Record<string, number>; 
  score: number;
}

// 5. 收藏状态统计
export interface BangumiCollection {
  on_hold: number;
  dropped: number;
  wish: number;
  collect: number;
  doing: number;
}

// 6. 核心：单个条目 (Subject)
export interface BangumiSubject {
  id: number;
  type: number; // 1:书籍, 2:动画, 3:音乐, 4:游戏, 6:三次元
  name: string;
  name_cn: string; // 可能为空字符串
  summary: string;
  date: string; // YYYY-MM-DD
  platform: string; // 如 "漫画", "TV", "小说"
  
  image: string; // 封面图 URL (通常等于 images.large)
  
  tags: BangumiTag[];
  infobox: InfoboxItem[]; // 有些条目可能没有 infobox
  
  rating: BangumiRating;
  collection: BangumiCollection;
  
  eps: number;
  total_episodes: number;
  volumes: number;
  series: boolean;
  locked: boolean;
  nsfw: boolean;
  meta_tags: string[];
}

// 7. 根对象：搜索响应结果
export interface BangumiSearchResponse {
  data: BangumiSubject[];
  total: number;
  limit: number;
  offset: number;
}
