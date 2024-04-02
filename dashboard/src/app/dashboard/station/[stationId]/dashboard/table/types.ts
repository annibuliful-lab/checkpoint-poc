import { DevicePermittedLabel } from "@/apollo-client";

export type ImsiImeiTransaction = {
  id: string;
  arrivalTime: Date;
  imei: { status: DevicePermittedLabel; imei: number };
  imsi: { status: DevicePermittedLabel; imsi: number };
  phoneModel: string;
  tags: string[];
};
