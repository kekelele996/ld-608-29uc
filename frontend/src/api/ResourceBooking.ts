import { request } from "./http";
import { mockBookings } from "../mocks/seedData";
import type { ResourceBooking } from "../types/ResourceBooking";

export async function listResourceBookings(): Promise<ResourceBooking[]> {
  try {
    return await request<ResourceBooking[]>("/resource-bookings");
  } catch {
    return structuredClone(mockBookings);
  }
}
