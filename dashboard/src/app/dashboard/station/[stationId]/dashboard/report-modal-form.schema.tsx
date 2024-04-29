import * as Yup from "yup";

// ----------------------------------------------------------------------

export const ReportDashnoardFormSchema = Yup.object().shape({
  note: Yup.string().required(),
});
