import {
  Button,
  Card,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Stack,
} from "@mui/material";
import React from "react";
import _ from "lodash";
import LoadingButton from "@mui/lab/LoadingButton";
import { ReportDashnoardFormSchema } from "./report-modal-form.schema";
import { useForm } from "react-hook-form";
import FormProvider from "@/components/hook-form/form-provider";
import { yupResolver } from "@hookform/resolvers/yup";
import { RHFTextField } from "@/components/hook-form";

type Props = {
  opened: boolean;
  onClose: () => void;
  //   stationId: string;
  //   row?: StationDashboardActivity;
};
export default function CreateReportModal({
  opened,
  onClose,
}: //   stationId,
//   row,
Props) {
  const title = "Report Issue";
  // IStationItem

  const methods = useForm({
    resolver: yupResolver(ReportDashnoardFormSchema),
  });

  const {
    reset,
    handleSubmit,
    formState: { isSubmitting },
  } = methods;
  const onSubmit = handleSubmit(async (values) => {
    try {
     
      const input = values;

      //   if (dafaultValues?.id) {
      //     await update({
      //       variables: { id: dafaultValues?.id, ...input },
      //       refetchQueries: [GetMobileDeviceConfigurationsDocument],
      //     });
      //     reset(input);
      //   } else {
      //     await create({
      //       variables: input,
      //       refetchQueries: [GetMobileDeviceConfigurationsDocument],
      //     });
      //     reset();
      //   }
      //   createResponse.reset();
      //   updateResponse.reset();
      //   onClose();
    } catch (error) {
      console.error(error);
    }
  });

  return (
    <Dialog fullWidth open={opened} onClose={onClose} maxWidth="sm">
      <DialogTitle sx={{ pb: 2 }}>{title}</DialogTitle>
      <DialogContent
        sx={{
          overflow: "unset",
        }}
      >
        <FormProvider methods={methods} onSubmit={onSubmit}>
          <Stack spacing={0.5}>
            <Stack direction={"row"} spacing={0.5}>
              <RHFTextField
                name="note"
                multiline
                rows={5}
                sx={{
                  width: 1,
                }}
                placeholder="report an issue here..."
              />
            </Stack>
            <DialogActions>
              <Button variant="outlined" onClick={onClose}>
                Cancel
              </Button>

              <LoadingButton
                type="submit"
                variant="contained"
                color="primary"
                loading={isSubmitting}
              >
                Submit
              </LoadingButton>
            </DialogActions>
          </Stack>
        </FormProvider>
      </DialogContent>
    </Dialog>
  );
}
