import type { APIRoute, GetStaticPaths } from 'astro';
import type { Doc } from '../../lib/docs';
import { llmsDocs, pageMarkdown } from '../../lib/llms';

// Each English page as Markdown at /docs/<id>.md, the first at /docs/index.md.
export const getStaticPaths: GetStaticPaths = async () => {
  const { docs } = await llmsDocs();
  return docs.map((doc) => ({ params: { slug: doc.id }, props: { doc } }));
};

export const GET: APIRoute = ({ props }) =>
  new Response(pageMarkdown((props as { doc: Doc }).doc), { headers: { 'Content-Type': 'text/markdown; charset=utf-8' } });
