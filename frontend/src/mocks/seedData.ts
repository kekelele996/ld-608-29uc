export const mockData = {
  "flightTurnaround": [
    {
      "id": 1,
      "flight_no": "CA1234",
      "aircraft_reg": "B-6688",
      "stand_no": "12",
      "arrival_time": "2026-10-07T04:50:00Z",
      "departure_time": "2026-10-07T08:30:00Z",
      "turnaround_status": "IN_SERVICE",
      "delay_reason": "加油车调度延迟"
    },
    {
      "id": 2,
      "flight_no": "MU5678",
      "aircraft_reg": "B-8866",
      "stand_no": "07",
      "arrival_time": "2026-10-07T06:10:00Z",
      "departure_time": "2026-10-07T09:40:00Z",
      "turnaround_status": "ON_STAND",
      "delay_reason": "行李分拣系统故障"
    },
    {
      "id": 3,
      "flight_no": "ZH9012",
      "aircraft_reg": "B-1234",
      "stand_no": "21",
      "arrival_time": "2026-10-07T08:20:00Z",
      "departure_time": "2026-10-07T11:00:00Z",
      "turnaround_status": "DELAYED",
      "delay_reason": "前序航班晚到"
    }
  ],
  "groundTask": [
    {
      "id": 1,
      "turnaround_id": 1,
      "task_type": "CLEANING",
      "team_id": 1,
      "planned_start": "2026-10-07T01:30:00Z",
      "deadline": "2026-10-07T02:00:00Z",
      "actual_finish": "",
      "status": "PENDING",
      "blocker_note": ""
    },
    {
      "id": 2,
      "turnaround_id": 1,
      "task_type": "BAGGAGE",
      "team_id": 2,
      "planned_start": "2026-10-07T05:40:00Z",
      "deadline": "2026-10-07T06:20:00Z",
      "actual_finish": "2026-10-07T06:55:00Z",
      "status": "SIGNED",
      "blocker_note": ""
    },
    {
      "id": 3,
      "turnaround_id": 1,
      "task_type": "CATERING",
      "team_id": 3,
      "planned_start": "2026-10-07T06:00:00Z",
      "deadline": "2026-10-07T06:50:00Z",
      "actual_finish": "2026-10-07T07:05:00Z",
      "status": "DONE",
      "blocker_note": ""
    },
    {
      "id": 4,
      "turnaround_id": 2,
      "task_type": "REFUEL",
      "team_id": 1,
      "planned_start": "2026-10-07T07:00:00Z",
      "deadline": "2026-10-07T08:00:00Z",
      "actual_finish": "",
      "status": "PENDING",
      "blocker_note": ""
    },
    {
      "id": 5,
      "turnaround_id": 2,
      "task_type": "WATER_SERVICE",
      "team_id": 2,
      "planned_start": "2026-10-07T06:30:00Z",
      "deadline": "2026-10-07T07:00:00Z",
      "actual_finish": "2026-10-07T07:20:00Z",
      "status": "SIGNED",
      "blocker_note": ""
    },
    {
      "id": 6,
      "turnaround_id": 3,
      "task_type": "PUSHBACK",
      "team_id": 3,
      "planned_start": "2026-10-07T09:30:00Z",
      "deadline": "2026-10-07T10:00:00Z",
      "actual_finish": "",
      "status": "BLOCKED",
      "blocker_note": "牵引车故障待修"
    }
  ],
  "groundResource": [
    {
      "id": 1,
      "resource_code": "resource code 1",
      "resource_type": "CATERING",
      "location": "location 1",
      "availability_status": "ON_STAND",
      "maintenance_due_at": "2026-06-11T09:00:00Z",
      "owner_team": "owner team 1"
    },
    {
      "id": 2,
      "resource_code": "resource code 2",
      "resource_type": "BAGGAGE",
      "location": "location 2",
      "availability_status": "IN_SERVICE",
      "maintenance_due_at": "2026-06-12T09:00:00Z",
      "owner_team": "owner team 2"
    },
    {
      "id": 3,
      "resource_code": "resource code 3",
      "resource_type": "REFUEL",
      "location": "location 3",
      "availability_status": "ARRIVING",
      "maintenance_due_at": "2026-06-13T09:00:00Z",
      "owner_team": "owner team 3"
    }
  ],
  "resourceBooking": [
    {
      "id": 1,
      "resource_id": 1,
      "turnaround_id": 1,
      "task_id": 1,
      "start_time": "2026-06-11T09:00:00Z",
      "end_time": "2026-06-11T09:00:00Z",
      "booking_status": "ON_STAND",
      "conflict_reason": "conflict reason 1"
    },
    {
      "id": 2,
      "resource_id": 2,
      "turnaround_id": 2,
      "task_id": 2,
      "start_time": "2026-06-12T09:00:00Z",
      "end_time": "2026-06-12T09:00:00Z",
      "booking_status": "IN_SERVICE",
      "conflict_reason": "conflict reason 2"
    },
    {
      "id": 3,
      "resource_id": 3,
      "turnaround_id": 3,
      "task_id": 3,
      "start_time": "2026-06-13T09:00:00Z",
      "end_time": "2026-06-13T09:00:00Z",
      "booking_status": "ARRIVING",
      "conflict_reason": "conflict reason 3"
    }
  ],
  "delayEvent": [
    {
      "id": 1,
      "turnaround_id": 1,
      "delay_type": "REFUEL",
      "minutes": 30,
      "root_cause": "加油车调度延迟",
      "responsibility_team": "机务一队",
      "resolved_at": ""
    },
    {
      "id": 2,
      "turnaround_id": 1,
      "delay_type": "CATERING",
      "minutes": 15,
      "root_cause": "配餐车晚到",
      "responsibility_team": "配餐组",
      "resolved_at": "2026-10-07T06:40:00Z"
    },
    {
      "id": 3,
      "turnaround_id": 2,
      "delay_type": "BAGGAGE",
      "minutes": 45,
      "root_cause": "行李分拣系统故障",
      "responsibility_team": "行李组",
      "resolved_at": ""
    }
  ]
} as const;
