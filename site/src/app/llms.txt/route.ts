import { buildLlmsTxt } from "@/lib/llms-txt";

export const dynamic = "force-static";

/** Well-known LLM index at /llms.txt — note-20260808-26ff4eb4. */
export function GET() {
  return new Response(buildLlmsTxt(), {
    headers: {
      "Content-Type": "text/plain; charset=utf-8",
    },
  });
}
