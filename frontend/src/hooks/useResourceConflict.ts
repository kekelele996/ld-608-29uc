import { useMemo } from "react";
import type { ResourceBooking } from "../types/ResourceBooking";

/** 从预约列表中筛出冲突预约，并给出冲突资源 id 集合（资源页/看板共用）。 */
export function useResourceConflict(bookings: ResourceBooking[]) {
  return useMemo(() => {
    const conflicts = bookings.filter((b) => b.booking_status === "CONFLICT");
    const conflictResourceIds = new Set(conflicts.map((b) => b.resource_id));
    return { conflicts, conflictResourceIds, conflictCount: conflicts.length };
  }, [bookings]);
}
