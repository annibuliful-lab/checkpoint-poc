import {
  TableEmptyRows,
  TableHeadCustom,
  TableNoData,
  TablePaginationCustom,
  TableSkeleton,
  useTable,
} from "@/components/table";
import { Card, Table, TableBody, TableContainer } from "@mui/material";
import React, { useCallback, useMemo, useState } from "react";
import Scrollbar from "@/components/scrollbar";
import VehicleTableRow from "./table-row";
import { TABLE_HEAD } from "./const";
import { useBoolean } from "@/hooks/use-boolean";
import _ from "lodash";
import { PropWithStationLocationId } from "../../types";
import { StationDashboardActivity, useGetStationDashboardActivityByIdLazyQuery, useGetStationDashboardActivityByIdQuery, useGetStationVehicleActivitiesQuery, useGetStationVehicleActivityByIdLazyQuery } from "@/apollo-client";
import { VehicleTransection, transformData } from "./types";
import VehicleInfoModal from "./vehicle-info-Modal";

const defaultFilters = {
  name: "",
};
export default function VehicleTransectionTable({
  stationLocationId,
}: PropWithStationLocationId) {

  const [getVehicleDashboardInfo] = useGetStationDashboardActivityByIdLazyQuery()
  const table = useTable({
    defaultOrderBy: "createdAt",
    defaultOrder: "desc",
    defaultRowsPerPage: 10,
  });
  const { data, loading } = useGetStationVehicleActivitiesQuery({
    variables: {
      limit: 0,
      skip: 0,
      stationId: stationLocationId,
    },
  });
  const [filters] = useState(defaultFilters);
  const [vehicleInfo,setVehicleInfo] = useState<StationDashboardActivity | undefined>(undefined)
  const dataInTable = useMemo(
    () =>
      data?.getStationVehicleActivities?.map((row) => transformData(row)) ?? [],
    [data?.getStationVehicleActivities]
  );
  const handleDeleteRow = useCallback((id: string) => {
    //
  }, []);
  const openVehicleForm = useBoolean();

  const openModalInfo = async (row : VehicleTransection) =>{
    
    const result = await getVehicleDashboardInfo({
      variables: {
        id: row.id
      }
    })
    setVehicleInfo(result.data?.getStationDashboardActivityById || undefined )
    openVehicleForm.onTrue()
  }

  return (
    <Card>
      <VehicleInfoModal
        row={vehicleInfo}
        opened={openVehicleForm.value}
        onClose={openVehicleForm.onFalse}
        stationId={stationLocationId}
      />
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
                      onClick={() => openModalInfo(row)}
                      // onDeleteRow={() => handleDeleteRow(row.id)}
                      // onEdit={() => console.log(row)}
                    />
                  ))}
                <TableEmptyRows
                  height={76}
                  emptyRows={dataInTable.length}
                  colSpan={TABLE_HEAD.length}
                />
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
