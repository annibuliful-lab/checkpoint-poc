import { DevicePermittedLabel } from "@/apollo-client";

export type VehicleTransection = {
  id: string;
  arrivalTime: Date;
  licensePlate: {
    image?: string;
    license: string;
    type: string;
    status: DevicePermittedLabel;
  };
  brand: string;
  vehicle: {
    type: string;
  };
  color: {
    name: string;
    code: string;
  };
  stationSite: string;
  imei: {
    list: string[];
    total: number;
  };
  imsi: {
    list: string[];
    total: number;
  };
  tags: string[];
  remark: string;
};
