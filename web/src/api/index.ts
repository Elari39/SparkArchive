import { apiGet, qs } from './client'
import {
  EventItem,
  GraphData,
  PersonDetail,
  PersonSummary,
  SearchResult,
  SiteMeta,
  TermDetail,
  TermSummary,
  WorkDetail,
  WorkSummary,
} from './types'
import { z } from 'zod'

export const getMeta = (signal?: AbortSignal) => apiGet('/meta', SiteMeta, signal)
export const getPeople = (signal?: AbortSignal) => apiGet('/people', z.array(PersonSummary), signal)
export const getPerson = (slug: string, signal?: AbortSignal) =>
  apiGet(`/people/${encodeURIComponent(slug)}`, PersonDetail, signal)
export const getWorks = (p: { person?: string; q?: string } = {}, signal?: AbortSignal) =>
  apiGet('/works' + qs(p), z.array(WorkSummary), signal)
export const getWork = (slug: string, signal?: AbortSignal) =>
  apiGet(`/works/${encodeURIComponent(slug)}`, WorkDetail, signal)
export const getEvents = (p: { person?: string; category?: string } = {}, signal?: AbortSignal) =>
  apiGet('/events' + qs(p), z.array(EventItem), signal)
export const getTerms = (p: { q?: string } = {}, signal?: AbortSignal) =>
  apiGet('/terms' + qs(p), z.array(TermSummary), signal)
export const getTerm = (slug: string, signal?: AbortSignal) =>
  apiGet(`/terms/${encodeURIComponent(slug)}`, TermDetail, signal)
export const getGraph = (signal?: AbortSignal) => apiGet('/relations', GraphData, signal)
export const search = (
  p: { q: string; type?: string; limit?: number; offset?: number },
  signal?: AbortSignal,
) => apiGet('/search' + qs(p), SearchResult, signal)
