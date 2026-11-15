import type { APIRoute } from 'astro';
import { llmsDocs, llmsFull } from '../lib/llms';

export const GET: APIRoute = async () => {
  const { docs } = await llmsDocs();
  return new Response(llmsFull(docs), { headers: { 'Content-Type': 'text/plain; charset=utf-8' } });
};
