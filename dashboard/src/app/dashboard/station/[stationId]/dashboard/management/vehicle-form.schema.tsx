import {
  BlacklistPriority,
  CreateVehicleTargetConfigurationMutationVariables,
  DevicePermittedLabel,
} from "@/apollo-client";
import * as Yup from "yup";

// ----------------------------------------------------------------------
export const vehicleFormDefaultValues: CreateVehicleTargetConfigurationMutationVariables =
  {
    prefix: "",
    number: "",
    province: "",
    type: "",
    country: "",
    permittedLabel: DevicePermittedLabel.None,
    blacklistPriority: BlacklistPriority.Normal,
    tags: [],
    stationLocationId: "",
    color: "",
    brand: "",
  };
export type VehicleFormInput = typeof vehicleFormDefaultValues;
export const VehicleFormSchema = Yup.object().shape({
  prefix: Yup.string().required(),
  number: Yup.string().required(),
  province: Yup.string().required(),
  type: Yup.string().required(),
  permittedLabel: Yup.string().required(),
  blacklistPriority: Yup.string().required(),
  country: Yup.string().notRequired(),
  tags: Yup.array().of(Yup.string()).max(5),
});
