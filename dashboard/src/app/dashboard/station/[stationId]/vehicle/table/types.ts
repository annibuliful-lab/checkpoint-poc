import {
  DevicePermittedLabel,
  StationVehicleActivity,
  StationVehicleActivityTag,
} from "@/apollo-client";

export type VehicleTransection = {
  id: string;
  arrivalTime: string;
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
  tags: StationVehicleActivityTag[];
  remark: string;
};
export function transformData(
  graphqlData: StationVehicleActivity
): VehicleTransection {
  const vehicleTransection: VehicleTransection = {
    id: graphqlData.id,
    arrivalTime: graphqlData.arrivalTime,
    licensePlate: {
      image: graphqlData.licensePlate.image ?? "",
      license: graphqlData.licensePlate.license,
      type: graphqlData.licensePlate.type,
      status: graphqlData.licensePlate.status,
    },
    brand: graphqlData.brand,
    vehicle: {
      type: graphqlData.vehicle.type,
    },
    color: {
      name: graphqlData.color.name,
      code: graphqlData.color.code,
    },
    stationSite: graphqlData.stationSite,
    imei: {
      list: graphqlData.imei.list ?? [],
      total: graphqlData.imei.total,
    },
    imsi: {
      list: graphqlData.imsi.list ?? [],
      total: graphqlData.imsi.total,
    },
    tags: graphqlData.tags ?? [],
    remark: graphqlData.remark,
  };

  return vehicleTransection;
}
