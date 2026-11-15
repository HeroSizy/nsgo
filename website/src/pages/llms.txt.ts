import type { APIRoute } from 'astro';
import { llmsDocs, llmsIndex } from '../lib/llms';

export const GET: APIRoute = async () => {
  const { groups } = await llmsDocs();
  return new Response(llmsIndex(groups), { headers: { 'Content-Type': 'text/plain; charset=utf-8' } });
};
