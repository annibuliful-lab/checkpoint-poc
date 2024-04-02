import * as Yup from "yup";

// ----------------------------------------------------------------------
export type ResponsiblePerson = {
  name: string;
  phone: string;
};
export const stationFormDefaultValues = {
  title: "",
  description: "",
  department: "",
  latitude: 13.787220202859372,
  longitude: 100.56834202403502,
  tags: [] as string[],
  remark: "",
  responsiblePersons: [] as ResponsiblePerson[],
};
export type StationFormInput = typeof stationFormDefaultValues;
export const StationFormSchema = Yup.object().shape({
  title: Yup.string().required(),
  description: Yup.string().notRequired(),
  department: Yup.string().required(),
  latitude: Yup.number().required(),
  longitude: Yup.number().required(),
  tags: Yup.array().of(Yup.string()).max(5),
  responsiblePersons: Yup.array()
    .of(
      Yup.object().shape({
        name: Yup.string().required(),
        phone: Yup.string().required(),
      })
    )
    .max(3),
});
