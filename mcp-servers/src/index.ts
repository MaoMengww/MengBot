import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { string, z } from "zod";
import {
    BangumiSubject,
    BangumiSubjectCalendar,
    CalendarResponse,
} from "./types.js";
import { ca } from "zod/v4/locales";
//初始化Server
const server = new McpServer({
    name: "Bangumi",
    version: "1.0.0",
})


server.registerTool(
    "serch_banguni_characters",
    {
        description: "搜索 Bangumi (番组计划) 角色信息，支持 NSFW 开关,可以通过属于角色名称获得一些角色的基本信息",
        inputSchema: z.object({
            keyword: z.string().describe("角色名称关键词，如 '阿尔托莉雅'"),
            filter: z.object({
                nsfw: z.boolean().optional().default(false).describe("是否包含NSFW内容"),
            }).optional().default({ nsfw: true }),
        }),
    },
    async ({ keyword, filter }) => {
        try {
            const apiUrl = `https://api.bgm.tv/v0/search/characters?limit=3`;
            const requestBody = {
                keyword,
                filter: {
                    nsfw: filter?.nsfw ?? false
                }
            };

            //发起请求
            const response = await fetch(apiUrl, {
                method: 'POST',
                headers: {
                    'User-Agent': 'MengBot/1.0.0 (https://github.com/maomeng/MengBot)',
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                },
                body: JSON.stringify(requestBody),
            });

            if (!response.ok) {
                return {
                    content: [{
                        type: "text",
                        text: `搜索角色 ${keyword} 失败，状态码 ${response.status}`
                    }],
                    isError: true,
                };
            }

            //解析响应
            const json = await response.json() as any;
            const list = json.data || [];

            if (!list || list.length === 0) {
                return {
                    content: [{
                        type: "text",
                        text: `搜索角色 ${keyword} 没有结果`
                    }]
                };
            }

            //取第一个角色
            const item = list[0] as BangumiSubject;

            //格式化infobox
            let infoboxTexT = "";
            if (item.infobox) {
                infoboxTexT = item.infobox
                    .map(i => {
                        const v = Array.isArray(i.value) ?
                            i.value.map(j => j.v).join("\n")
                            : i.value;
                        return `-${i.key}:${v} `;
                    })
                    .join("\n");
            }
            const resultTxT = `
                角色名称:${item.name}
                角色简介:${item.summary || "无"}
                基本信息:${infoboxTexT || "无"}
                性别:${item.gender || "未知"}
                评论数:${item.stat?.comments || 0}
                收藏数:${item.stat?.collects || 0}
            `.trim();

            return {
                content: [{
                    type: "text",
                    text: resultTxT
                }]
            };

        } catch (error) {
            return {
                content: [{
                    type: "text",
                    text: `搜索角色 ${keyword} 失败，错误信息 ${error}`
                }],
                isError: true,
            };
        }
    }
);

server.registerTool(
    "bangumi_calendar",
    {
        description: "根据输入的目标日期的星期几,获取当天会播放的番剧",
        inputSchema: z.object({
            day: z.string().describe("应输入星期一,星期二,星期三,星期四,星期五,星期六,星期日, 其它任何相近的意思均转为七个选项之一",),
        }),
    },
    async ({ day }) => {
        try {
            const apiUrl = `https://api.bgm.tv/calendar`;
            //发起请求
            const response = await fetch(apiUrl, {
                method: 'GET',
                headers: {
                    'User-Agent': 'MengBot/1.0.0 (https://github.com/maomeng/MengBot)',
                },
            });
            if (!response.ok) {
                return {
                    content: [{
                        type: "text",
                        text: `获取 ${day} 番剧列表失败，状态码 ${response.status}`
                    }],
                    isError: true,
                };
            }

            //解析响应
            const json = await response.json() as CalendarResponse;
            const targetDay = day ;

            const dayData = json.find(d => d.weekday.cn === targetDay);
            if (!dayData) {
                return {
                    content: [{
                        type: "text",
                        text: `获取 ${day} 番剧列表失败，没有数据`
                    }],
                    isError: true,
                };
            }

            //格式化数据
            const result = dayData.items.map(item => ({
                name: item.name,
                name_cn: item.name_cn || null,
                score: item.rating?.score || null,
                rank: item.rank || null,
                air_date: item.air_date || null,
            }));

            return {
                content: [{
                    type: "text",
                    text: `
                    获取 ${day} 番剧列表成功，共 ${result.length} 个番剧
                    ${result.map(i => `${i.name_cn || i.name} (${i.score || "无"}分, aired on ${i.air_date || "无"})`).join("\n")}
                    `
                }],
            };
            
        }
        catch (error) {
            return {
                content: [{
                    type: "text",
                    text: `获取 ${day} 番剧列表失败，错误信息 ${error}`
                }],
                isError: true,
            };
        }
    }
);

server.registerTool(
    "bangumi_search",
    {
        description: `
        此工具通过调用bangumi的api来搜索书籍，动漫(或者电视剧，电影，综艺等)，游戏，和音乐的详细信息,
        `,
        inputSchema: z.object({
            keyword: z.string().describe("书籍名称，动漫名称，游戏名称，或音乐名称"),
            filter: z.object({
                type: z.array(z.number()).describe("搜索类型，1为书籍，2为动漫(或者电视剧，电影，综艺等)，3为游戏，4为音乐, 可以多选"),
                tag: z.array(z.string()).describe("搜索标签，例如 童年，原创等(非必要,建议为空,不填)").nullable(),
            })
        }),
    },
    async ({ keyword, filter }) => {
        try {
            const apiUrl = `https://api.bgm.tv/v0/`;
            const requestBody = {
                keyword: keyword,
                filter: {
                    type: filter.type.join(",") || [2],
                    tag: filter.tag?.join(",") || null,
                }
            }
            //发起请求
            const response = await fetch(apiUrl, {
                method: 'POST',
                headers: {
                    'User-Agent': 'MengBot/1.0.0 (https://github.com/MaoMengww/MengBot)',
                },
                body: JSON.stringify(requestBody),
            });
            if (!response.ok) {
                return {
                    content: [{
                        type: "text",
                        text: `搜索 ${keyword} 失败，状态码 ${response.status}`
                    }],
                    isError: true,
                };
            }

            const json = await response.json() as any;
            const list = json.data || []; 
            if (list.length === 0) {
                return {
                    content: [{
                        type: "text",
                        text: `搜索 ${keyword} 失败，没有数据`
                    }],
                    isError: true,
                };
            }

            //格式化数据
            const result = list[0] as BangumiSubject;

            let infoboxTexT = "";
            if (result.infobox) {
                infoboxTexT = result.infobox
                    .map(i => {
                        const v = Array.isArray(i.value) ?
                            i.value.map(j => j.v).join("\n")
                            : i.value;
                        return `-${i.key}:${v} `;
                    })
                    .join("\n");
            }

            let tagText = "";
            if (result.tags) {
                tagText = result.tags.map(i => i.name).join(",");
            }

            const typeText = result.type === 2 ? "动漫" : result.type === 3 ? "音乐" : result.type === 4 ? "游戏" : "书籍";

            const resultText = `
                名称:${result.name_cn || result.name}
                简介:${result.summary || "无"}
                发行时间：${result.date || "无"}
                类型:${typeText || "无"}
                platform:${result.platform || "无"}
                标签:${tagText || "无"}
                评分:${result.rating?.score || "无"}分
                章节数:${result.eps || "无"}
                总集数:${result.total_episodes || "无"}
                卷数:${result.volumes || "无"}
                是系列:${result.series ? "是" : "否"}
                更多:${infoboxTexT}
            `

            return {
                content: [{
                    type: "text",
                    text: resultText
                }],
            };
        } catch (error) {
            return {
                content: [{
                    type: "text",
                    text: `搜索 ${keyword} 失败，错误信息 ${error}`
                }],
                isError: true,
            };
        }
    }
);

async function main() {
    const transport = new StdioServerTransport();
    await server.connect(transport);
    console.error("Bangumi server started");
}

main().catch(console.error);
