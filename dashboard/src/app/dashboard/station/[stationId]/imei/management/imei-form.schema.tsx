import {
  BlacklistPriority,
  CreateMobileDeviceConfigurationMutationVariables,
  DevicePermittedLabel,
} from "@/apollo-client";
import * as Yup from "yup";

// ----------------------------------------------------------------------
export const imeiFormDefaultValues: CreateMobileDeviceConfigurationMutationVariables =
  {
    imei: "",
    imsi: "",
    tags: [],
    stationLocationId: "",
    title: "",
    permittedLabel: DevicePermittedLabel.None,
    blacklistPriority: BlacklistPriority.Normal,
  };
export type ImeiFormInput = typeof imeiFormDefaultValues;
export const ImeiFormSchema = Yup.object().shape({
  imei: Yup.string().notRequired(),
  imsi: Yup.string().notRequired(),
  permittedLabel: Yup.string().required(),
  blacklistPriority: Yup.string().required(),
  tags: Yup.array().of(Yup.string()).max(5),
}) as Yup.ObjectSchema<ImeiFormInput>;
