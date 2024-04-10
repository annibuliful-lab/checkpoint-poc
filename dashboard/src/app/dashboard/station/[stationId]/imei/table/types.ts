import {
  DevicePermittedLabel,
  StationImeiImsiActivity,
  StationImeiImsiActivityTag,
} from "@/apollo-client";

export type ImeiImsiTransaction = {
  id: string;
  arrivalTime: string;
  imei: string;
  imsi: string;
  phoneModel: string;
  licensePlate: string;
  stationSite: string;
  tags: StationImeiImsiActivityTag[];
  status: DevicePermittedLabel;
};

export function transformData(
  data: StationImeiImsiActivity
): ImeiImsiTransaction {
  const vehicleTransection: ImeiImsiTransaction = {
    id: data.id,
    arrivalTime: data.arrivalTime,
    imei: data.imei?.imei ?? "",
    imsi: data.imsi?.imsi ?? "",
    phoneModel: data.phoneModel ?? "",
    licensePlate: data.licensePlate ?? "",
    stationSite: data.stationSiteName ?? "",
    tags: data.tags ?? [],
    status:
      data.imei?.permittedLabel ||
      data.imsi?.permittedLabel ||
      DevicePermittedLabel.None,
  };

  return vehicleTransection;
}
