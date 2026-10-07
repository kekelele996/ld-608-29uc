import type { GroundResource } from "../types/GroundResource";

export const createDefaultGroundResource = (overrides: Partial<GroundResource> = {}): GroundResource => ({
  id: 0,
  resource_code: "",
  resource_type: "POWER",
  location: "",
  availability_status: "AVAILABLE",
  status_text: "可用",
  maintenance_due_at: null,
  owner_team: "",
  ...overrides
});

export const createGroundResourceForm = createDefaultGroundResource;
export const createGroundResourceResponse = createDefaultGroundResource;
