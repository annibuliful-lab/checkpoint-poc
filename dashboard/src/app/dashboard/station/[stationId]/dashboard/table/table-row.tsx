import Stack from "@mui/material/Stack";
import Button from "@mui/material/Button";
import MenuItem from "@mui/material/MenuItem";
import TableRow from "@mui/material/TableRow";

import TableCell from "@mui/material/TableCell";
import IconButton from "@mui/material/IconButton";

import { useBoolean } from "@/hooks/use-boolean";
import Iconify from "@/components/iconify";

import CustomPopover, { usePopover } from "@/components/custom-popover";
import Label from "@/components/label";
import { ConfirmDialog } from "@/components/custom-dialog";
import { Edit } from "@mui/icons-material";
import { Alert, Typography } from "@mui/material";
import { useRouter } from "@/routes/hooks";
import { fDateTime } from "@/utils/format-time";
import { DevicePermittedLabel } from "@/apollo-client";
import { ImsiImeiTransaction } from "./types";

// ----------------------------------------------------------------------

type Props = {
  row: ImsiImeiTransaction;
  onDeleteRow: VoidFunction;
  onEdit: VoidFunction;
};

export default function VehicleTableRow({ row, onDeleteRow, onEdit }: Props) {
  const { id, tags, arrivalTime, imei, imsi, phoneModel } = row;

  const isBlackList = [imei.status, imsi.status].includes(
    DevicePermittedLabel.Blacklist
  );
  const isWhiteList = [imei.status, imsi.status].includes(
    DevicePermittedLabel.Whitelist
  );
  return (
    <>
      <TableRow
        sx={{
          bgcolor:
            (isBlackList && "error.dark") ||
            (isWhiteList && "success.dark") ||
            undefined,
        }}
      >
        <TableCell>
          <Typography
            color={isBlackList || isWhiteList ? "white" : undefined}
            fontSize={14}
          >
            {fDateTime(arrivalTime)}
          </Typography>
        </TableCell>
        <TableCell>
          <Typography
            color={isBlackList || isWhiteList ? "white" : undefined}
            fontSize={14}
          >
            {imei.imei}
          </Typography>
        </TableCell>
        <TableCell>
          <Typography
            color={isBlackList || isWhiteList ? "white" : undefined}
            fontSize={14}
          >
            {imsi.imsi}
          </Typography>
        </TableCell>
        <TableCell>
          <Typography
            color={isBlackList || isWhiteList ? "white" : undefined}
            fontSize={14}
          >
            {phoneModel}
          </Typography>
        </TableCell>
        <TableCell>
          <Stack
            alignItems={"center"}
            spacing={1}
            direction={"row"}
            flexWrap={"wrap"}
          >
            {tags?.map((tag) => (
              <Label
                key={tag}
                variant={isBlackList || isWhiteList ? "filled" : "soft"}
                color="warning"
              >
                {tag}
              </Label>
            ))}
          </Stack>
        </TableCell>
      </TableRow>
    </>
  );
}
