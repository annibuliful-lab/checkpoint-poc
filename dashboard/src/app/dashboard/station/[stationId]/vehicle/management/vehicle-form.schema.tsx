import {
  BlacklistPriority,
  CreateVehicleTargetConfigurationMutationVariables,
  DevicePermittedLabel,
} from "@/apollo-client";
import * as Yup from "yup";

// ----------------------------------------------------------------------
export const vehicleFormDefaultValues: CreateVehicleTargetConfigurationMutationVariables & {
  vehicleBrand: "";
  vehicleColorName: "";
  vehicleType: "";
} = {
  prefix: "",
  number: "",
  province: "",
  type: "",
  country: "",
  vehicleBrand: "",
  vehicleColorName: "",
  vehicleType: "",
  permittedLabel: DevicePermittedLabel.None,
  blacklistPriority: BlacklistPriority.Normal,
  tags: [],
};
export type VehicleFormInput = typeof vehicleFormDefaultValues;
export const VehicleFormSchema = Yup.object().shape({
  prefix: Yup.string().required(),
  number: Yup.string().required(),
  province: Yup.string().required(),
  vehicleBrand: Yup.string().required(),
  vehicleColorName: Yup.string().required(),
  vehicleType: Yup.string().required(),
  type: Yup.string().required(),
  permittedLabel: Yup.string().required(),
  blacklistPriority: Yup.string().required(),
  country: Yup.string().notRequired(),
  tags: Yup.array().of(Yup.string()).max(5),
});
