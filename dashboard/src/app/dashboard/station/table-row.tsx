import Stack from "@mui/material/Stack";
import Button from "@mui/material/Button";
import MenuItem from "@mui/material/MenuItem";
import TableRow from "@mui/material/TableRow";

import TableCell from "@mui/material/TableCell";
import IconButton from "@mui/material/IconButton";

import { useBoolean } from "@/hooks/use-boolean";
import Iconify from "@/components/iconify";

import { IStationItem } from "./types";
import CustomPopover, { usePopover } from "@/components/custom-popover";
import Label from "@/components/label";
import { ConfirmDialog } from "@/components/custom-dialog";
import { Edit } from "@mui/icons-material";
import { Alert, Typography } from "@mui/material";
import { useRouter } from "@/routes/hooks";
import { paths } from "@/routes/paths";
import { useCallback } from "react";
import { useAtom, useAtomValue, useSetAtom } from "jotai";
import { stationHistoriesAtom } from "./store";

// ----------------------------------------------------------------------

type Props = {
  row: IStationItem;
  onDeleteRow: VoidFunction;
  onEdit: VoidFunction;
};

export default function StationTableRow({ row, onDeleteRow, onEdit }: Props) {
  const { id, latitude, longitude, remark, tags, title } = row;
  const confirm = useBoolean();
  const popover = usePopover();
  const route = useRouter();
  const setStationHistories = useSetAtom(stationHistoriesAtom);
  const onJoin = useCallback(() => {
    setStationHistories((prevs) => ({ ...prevs, [id]: row }));
    route.push(paths.dashboard.station.join(id).dashboard.root);
  }, [id, route, row, setStationHistories]);
  return (
    <>
      <TableRow hover>
        <TableCell>
          <Typography color={"text.primary"} onClick={onJoin}>
            {title}
          </Typography>
        </TableCell>
        <TableCell>
          <Label
            variant="soft"
            // color={
            //   (status === "online" && "success") ||
            //   (status === "waring" && "warning") ||
            //   (status === "offline" && "error") ||
            //   "default"
            // }
          >
            {/* {status} */}
            status
          </Label>
        </TableCell>
        <TableCell> - </TableCell>
        <TableCell>
          {latitude}, {longitude}
        </TableCell>
        <TableCell>
          <Stack
            alignItems={"center"}
            spacing={1}
            direction={"row"}
            flexWrap={"wrap"}
          >
            {tags?.map((tag) => (
              <Label key={tag.id}>{tag.title} </Label>
            ))}
          </Stack>
        </TableCell>
        <TableCell>0</TableCell>
        <TableCell>0</TableCell>
        <TableCell>0</TableCell>
        <TableCell>
          {remark && (
            <Alert
              variant="outlined"
              severity="warning"
              style={{ scale: 0.8, width: 200 }}
            >
              {remark}
            </Alert>
          )}
        </TableCell>
        <TableCell>
          <Stack alignItems={"center"} justifyContent={"center"}>
            <IconButton size="small" onClick={onEdit}>
              <Edit fontSize="small" />
            </IconButton>
          </Stack>
        </TableCell>
      </TableRow>
      <CustomPopover
        open={popover.open}
        onClose={popover.onClose}
        arrow="right-top"
        sx={{ width: 140 }}
      >
        <MenuItem
          onClick={() => {
            confirm.onTrue();
            popover.onClose();
          }}
          sx={{ color: "error.main" }}
        >
          <Iconify icon="solar:trash-bin-trash-bold" />
          Delete
        </MenuItem>
      </CustomPopover>

      <ConfirmDialog
        open={confirm.value}
        onClose={confirm.onFalse}
        title="Delete"
        content="Are you sure want to delete?"
        action={
          <Button variant="contained" color="error" onClick={onDeleteRow}>
            Delete
          </Button>
        }
      />
    </>
  );
}
