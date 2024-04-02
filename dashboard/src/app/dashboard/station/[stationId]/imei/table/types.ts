import { DevicePermittedLabel } from "@/apollo-client";

export type ImeiImsiTransaction = {
  id: string;
  arrivalTime: Date;
  imei: number;
  imsi: number;
  phoneModel: string;
  licensePlate: string;
  stationSite: string;
  tags: string[];
  status: DevicePermittedLabel;
};
