import { z } from 'zod'

export const PersonSummary = z.object({
  slug: z.string(),
  name: z.string(),
  name_en: z.string(),
  birth: z.string(),
  death: z.string(),
  birth_place: z.string(),
  nationality: z.string(),
  epithet: z.string(),
  portrait_key: z.string(),
  summary: z.string(),
  thesis: z.string(),
})
export type PersonSummary = z.infer<typeof PersonSummary>

export const PersonSection = z.object({
  kind: z.string(),
  title: z.string(),
  body_md: z.string(),
  ord: z.number(),
})

export const PersonTimelineItem = z.object({
  year: z.number(),
  date_text: z.string(),
  title: z.string(),
  body: z.string(),
})

export const PersonDetail = PersonSummary.extend({
  timeline: z.array(PersonTimelineItem),
  sections: z.array(PersonSection),
  works: z.array(z.object({ slug: z.string(), title: z.string(), year: z.number() })),
  terms: z.array(z.object({ slug: z.string(), term: z.string() })),
})
export type PersonDetail = z.infer<typeof PersonDetail>

export const WorkSummary = z.object({
  slug: z.string(),
  person: z.string(),
  person_name: z.string(),
  title: z.string(),
  year: z.number(),
  category: z.string(),
  summary: z.string(),
})
export type WorkSummary = z.infer<typeof WorkSummary>

export const WorkDetail = WorkSummary.extend({
  source: z.string(),
  source_url: z.string(),
  excerpts: z.array(z.object({ text: z.string(), note: z.string() })),
})
export type WorkDetail = z.infer<typeof WorkDetail>

export const EventItem = z.object({
  id: z.number(),
  year: z.number(),
  date_text: z.string(),
  title: z.string(),
  body: z.string(),
  location: z.string(),
  category: z.string(),
  people: z.array(z.object({ slug: z.string(), name: z.string() })),
})
export type EventItem = z.infer<typeof EventItem>

export const TermSummary = z.object({
  slug: z.string(),
  term: z.string(),
  aliases: z.array(z.string()),
  definition: z.string(),
})
export type TermSummary = z.infer<typeof TermSummary>

export const TermDetail = TermSummary.extend({
  body: z.string(),
  people: z.array(z.object({ slug: z.string(), name: z.string() })),
  related: z.array(z.object({ slug: z.string(), term: z.string() })),
  sources: z.array(z.object({ slug: z.string(), title: z.string() })),
})
export type TermDetail = z.infer<typeof TermDetail>

export const GraphData = z.object({
  nodes: z.array(z.object({ slug: z.string(), name: z.string(), epithet: z.string() })),
  edges: z.array(
    z.object({
      from: z.string(),
      to: z.string(),
      kind: z.string(),
      label: z.string(),
      note: z.string(),
    }),
  ),
})
export type GraphData = z.infer<typeof GraphData>

export const SearchHit = z.object({
  kind: z.string(),
  ref: z.string(),
  title: z.string(),
  snippet: z.string(),
  person: z.string(),
})
export type SearchHit = z.infer<typeof SearchHit>

export const SearchResult = z.object({
  total: z.number(),
  items: z.array(SearchHit),
})
export type SearchResult = z.infer<typeof SearchResult>

export const SiteMeta = z.object({
  site: z.object({
    name: z.string(),
    name_en: z.string(),
    tagline: z.string(),
    subtitle: z.string(),
    description: z.string(),
    footer_note: z.string(),
    license_note: z.string(),
  }),
  counts: z.object({
    people: z.number(),
    works: z.number(),
    events: z.number(),
    terms: z.number(),
  }),
})
export type SiteMeta = z.infer<typeof SiteMeta>
