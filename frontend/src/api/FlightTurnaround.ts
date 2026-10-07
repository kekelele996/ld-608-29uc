import { request } from "./http";
import { mockFlights } from "../mocks/seedData";
import type { FlightTurnaround } from "../types/FlightTurnaround";

export async function listFlightTurnarounds(): Promise<FlightTurnaround[]> {
  try {
    return await request<FlightTurnaround[]>("/flight-turnarounds");
  } catch {
    return structuredClone(mockFlights);
  }
}
