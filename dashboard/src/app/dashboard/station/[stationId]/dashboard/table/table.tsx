import {
  TableEmptyRows,
  TableHeadCustom,
  TableNoData,
  TablePaginationCustom,
  TableSkeleton,
  useTable,
} from "@/components/table";
import { Card, Table, TableBody, TableContainer } from "@mui/material";
import React, { useCallback, useEffect, useMemo, useState } from "react";
import Scrollbar from "@/components/scrollbar";
import VehicleTableRow from "./table-row";
import { TABLE_HEAD } from "./const";
import { useBoolean } from "@/hooks/use-boolean";
import _ from "lodash";
import { useGetStationDashboardActivitiesQuery } from "@/apollo-client";
import { PropWithStationLocationId } from "../../types";
import { useSetAtom } from "jotai";
import { stationDashboardActivityAtom } from "../store";

const defaultFilters = {
  name: "",
};
export default function VehicleTransectionTable({
  stationLocationId,
}: PropWithStationLocationId) {
  const table = useTable({
    defaultOrderBy: "createdAt",
    defaultOrder: "desc",
    defaultRowsPerPage: 10,
  });
  const [filters, setFilters] = useState(defaultFilters);
  const { data, loading } = useGetStationDashboardActivitiesQuery({
    variables: {
      limit: 1000,
      skip: 0,
      stationId: stationLocationId,
    },
  });
  const dataInTable = useMemo(
    () => data?.getStationDashboardActivities ?? [],
    [data?.getStationDashboardActivities]
  );
  const handleDeleteRow = useCallback((id: string) => {
  }, []);
  const openVehicleForm = useBoolean();
  const [editVehicle, setEditVehicle] = useState(undefined);
  const handleEdit = useCallback(
    (item: any) => {
      setEditVehicle(item);
      openVehicleForm.onTrue();
    },
    [openVehicleForm]
  );
  const setDashboardActivity = useSetAtom(stationDashboardActivityAtom);
  useEffect(() => {
    const recent = data?.getStationDashboardActivities?.[0];
    if (recent) {
      setDashboardActivity(recent);
    }
    return () => {
      //
    };
  }, [data?.getStationDashboardActivities, setDashboardActivity]);

  return (
    <Card>
      <TableContainer sx={{ position: "relative", overflow: "unset" }}>
        <Scrollbar>
          <Table size={"small"} sx={{ minWidth: 960 }}>
            <TableHeadCustom
              order={table.order}
              orderBy={table.orderBy}
              headLabel={TABLE_HEAD}
              rowCount={0}
              numSelected={table.selected.length}
              onSort={table.onSort}
            />
            {loading ? (
              <TableSkeleton />
            ) : (
              <TableBody>
                {dataInTable
                  .slice(
                    table.page * table.rowsPerPage,
                    table.page * table.rowsPerPage + table.rowsPerPage
                  )
                  .map((row) => (
                    <VehicleTableRow
                      key={row.id}
                      row={row}
                      onDeleteRow={() => handleDeleteRow(row.id)}
                      onEdit={() => handleEdit(row)}
                    />
                  ))}
                <TableEmptyRows height={76} emptyRows={dataInTable.length} />
                <TableNoData
                  notFound={
                    !!Object.values(filters).filter((a) => a).length &&
                    dataInTable.length === 0
                  }
                />
              </TableBody>
            )}
          </Table>
        </Scrollbar>
      </TableContainer>

      <TablePaginationCustom
        disabled={false}
        count={dataInTable.length}
        page={table.page}
        rowsPerPage={table.rowsPerPage}
        onPageChange={table.onChangePage}
        onRowsPerPageChange={table.onChangeRowsPerPage}
        //
        dense={table.dense}
        onChangeDense={table.onChangeDense}
      />
    </Card>
  );
}
