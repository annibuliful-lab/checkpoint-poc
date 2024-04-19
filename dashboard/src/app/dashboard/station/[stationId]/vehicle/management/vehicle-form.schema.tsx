import {
  BlacklistPriority,
  CreateVehicleTargetConfigurationMutationVariables,
  // CreateVehicleTargetConfigurationMutationVariables,
  DevicePermittedLabel,
  VehicleTargetConfiguration,
  VehicleTargetConfigurationImage
} from "@/apollo-client";
import * as Yup from "yup"; 

// // ----------------------------------------------------------------------
export const vehicleFormDefaultValues: CreateVehicleTargetConfigurationMutationVariables = {
  stationLocationId: "",
  color: "",
  brand: "",
  prefix: "",
  number: "",
  province: "",
  type: "",
  permittedLabel: DevicePermittedLabel.Blacklist,
  blacklistPriority: BlacklistPriority.Danger
};



export type VehicleFormInput = typeof vehicleFormDefaultValues;
export const VehicleFormSchema = Yup.object().shape({
  prefix: Yup.string().required(),
  stationLocationId: Yup.string().required(),
  number: Yup.string().required(),
  province: Yup.string().required(),
  id: Yup.string().notRequired(),
  color: Yup.string().required(),
  type: Yup.string().required(),
  permittedLabel: Yup.string().required(),
  brand: Yup.string().required(),
  blacklistPriority: Yup.string().required(),
  country: Yup.string().notRequired(),
  tags: Yup.array().of(Yup.string()).max(5).notRequired(),
  images : Yup.array()
  .of(
    Yup.object().shape({
      id: Yup.string().notRequired(),
      s3Key: Yup.string().required(),
      type: Yup.string().required(),
    })
  )
  .max(4),
});
