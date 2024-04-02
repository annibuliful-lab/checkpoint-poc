import { DevicePermittedLabel } from "@/apollo-client";
import { ImeiImsiTransaction } from "./types";

export const TABLE_HEAD = [
  { id: "title", label: "Arrival Time", width: 116 },
  { id: "title", label: "IMEI Number", width: 116 },
  { id: "title", label: "IMSI Number", width: 116 },
  { id: "title", label: "Phone model", width: 116 },
  { id: "title", label: "License Plate", width: 116 },
  { id: "title", label: "Station Site", width: 116 },
  { id: "title", label: "Tags", width: 116 },
];
export const IMEI_IMSI_TRANSECTIONS: ImeiImsiTransaction[] = [
  {
    id: "1",
    arrivalTime: new Date("2024-03-27T08:00:00Z"),
    imei: 123456789012345,
    imsi: 987654321098765,
    phoneModel: "ExamplePhone1",
    licensePlate: "ABC123",
    stationSite: "StationA",
    tags: ["tag1", "tag2"],
    status: DevicePermittedLabel.Blacklist,
  },
  {
    id: "2",
    arrivalTime: new Date("2024-03-27T09:30:00Z"),
    imei: 234567890123456,
    imsi: 876543210987654,
    phoneModel: "ExamplePhone2",
    licensePlate: "DEF456",
    stationSite: "StationB",
    tags: ["tag2", "tag3"],
    status: DevicePermittedLabel.Whitelist,
  },
  {
    id: "2",
    arrivalTime: new Date("2024-03-27T11:15:00Z"),
    imei: 345678901234567,
    imsi: 765432109876543,
    phoneModel: "ExamplePhone3",
    licensePlate: "GHI789",
    stationSite: "StationC",
    tags: ["tag3", "tag4"],
    status: DevicePermittedLabel.None,
  },
];
