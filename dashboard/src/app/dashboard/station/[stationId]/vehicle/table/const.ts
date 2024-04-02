import { DevicePermittedLabel } from "@/apollo-client";
import { VehicleTransection } from "./types";

export const TABLE_HEAD = [
  { id: "title", label: "Arrival Time", minWidth: 160 },
  { id: "title", label: "License Plate", minWidth: 140 },
  { id: "title", label: "License Plate Type", minWidth: 200 },
  { id: "title", label: "Brand", minWidth: 116 },
  { id: "title", label: "Vehicle Type", minWidth: 140 },
  { id: "title", label: "Color", minWidth: 116 },
  { id: "title", label: "Station Site", minWidth: 140 },
  { id: "title", label: "IMEI amount", minWidth: 140, align: "center" },
  { id: "title", label: "IMSI amount ", minWidth: 140, align: "center" },
  { id: "title", label: "Tags", minWidth: 160 },
  { id: "title", label: "Remark", minWidth: 160 },
];
export const VEHICLR_TRANSECTIONS: VehicleTransection[] = [
  {
    id: "1",
    arrivalTime: new Date("2024-03-27T08:00:00Z"),
    licensePlate: {
      license: "ABC123",
      type: "Standard",
      status: DevicePermittedLabel.None,
    },
    brand: "Toyota",
    vehicle: {
      type: "Sedan",
    },
    color: {
      name: "Red",
      code: "#FF0000",
    },
    stationSite: "Station A",
    imei: {
      list: ["861536030196001"],
      total: 15,
    },
    imsi: {
      list: ["520031234567890"],
      total: 10,
    },
    tags: ["Tag1", "Tag2"],
    remark: "No special remarks",
  },
  {
    id: "2",
    arrivalTime: new Date("2024-03-27T09:30:00Z"),
    licensePlate: {
      license: "XYZ789",
      type: "Standard",
      status: DevicePermittedLabel.Blacklist,
    },
    brand: "Honda",
    vehicle: {
      type: "SUV",
    },
    color: {
      name: "Blue",
      code: "#0000FF",
    },
    stationSite: "Station B",
    imei: {
      list: ["861536030196001"],
      total: 15,
    },
    imsi: {
      list: ["520031234567890"],
      total: 10,
    },
    tags: ["Tag3", "Tag4"],
    remark: "Needs maintenance soon",
  },
  {
    id: "3",
    arrivalTime: new Date("2024-03-27T09:30:00Z"),
    licensePlate: {
      license: "XYZ789",
      type: "Standard",
      status: DevicePermittedLabel.Whitelist,
    },
    brand: "Honda",
    vehicle: {
      type: "SUV",
    },
    color: {
      name: "Blue",
      code: "#0000FF",
    },
    stationSite: "Station B",
    imei: {
      list: ["1223434"],
      total: 10,
    },
    imsi: {
      list: ["xsw333445"],
      total: 15,
    },
    tags: ["Tag3", "Tag4"],
    remark: "Needs maintenance soon",
  },
];
