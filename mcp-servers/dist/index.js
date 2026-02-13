import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { z } from "zod";
//初始化Server
const server = new McpServer({
    name: "Bangumi",
    version: "1.0.0",
});
server.registerTool("serch_banguni_characters", {
    description: "搜索 Bangumi (番组计划) 角色信息，支持 NSFW 开关,可以通过属于角色名称获得一些角色的基本信息",
    inputSchema: z.object({
        keyword: z.string().describe("角色名称关键词，如 '阿尔托莉雅'"),
        filter: z.object({
            nsfw: z.boolean().optional().default(false).describe("是否包含NSFW内容"),
        }).optional().default({ nsfw: true }),
    }),
}, async ({ keyword, filter }) => {
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
        const json = await response.json();
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
        const item = list[0];
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
        const coverImage = item.images?.large || item.images?.medium || item.images?.grid || "";
        const resultTxT = `
                角色名称:${item.name}
                角色简介:${item.summary || "无"}
                基本信息:${infoboxTexT || "无"}
                性别:${item.gender || "未知"}
                评论数:${item.stat?.comments || 0}
                收藏数:${item.stat?.collects || 0}
                封面图片:${coverImage || "无"} 
            `.trim();
        return {
            content: [{
                    type: "text",
                    text: resultTxT
                }]
        };
    }
    catch (error) {
        return {
            content: [{
                    type: "text",
                    text: `搜索角色 ${keyword} 失败，错误信息 ${error}`
                }],
            isError: true,
        };
    }
});
async function main() {
    const transport = new StdioServerTransport();
    await server.connect(transport);
    console.error("Bangumi server started");
}
main().catch(console.error);
