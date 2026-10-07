import { request } from "./http";
import { mockResources } from "../mocks/seedData";
import type { GroundResource } from "../types/GroundResource";

export async function listGroundResources(): Promise<GroundResource[]> {
  try {
    return await request<GroundResource[]>("/ground-resources");
  } catch {
    return structuredClone(mockResources);
  }
}
