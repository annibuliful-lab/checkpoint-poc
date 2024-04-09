import {
  TableEmptyRows,
  TableHeadCustom,
  TableNoData,
  TablePaginationCustom,
  TableSkeleton,
  emptyRows,
  useTable,
} from "@/components/table";
import { Card, Table, TableBody, TableContainer } from "@mui/material";
import React, { useCallback, useMemo, useState } from "react";
import { IVehicleTableFilterValue, IVehicleTableFilters } from "./types";
import Scrollbar from "@/components/scrollbar";
import VehicleTableRow from "./table-row";
import { TABLE_HEAD } from "./const";
import { useBoolean } from "@/hooks/use-boolean";
import _ from "lodash";
import {
  GetVehicleTargetConfigurationsDocument,
  useDeleteVehicleTargetConfigurationMutation,
  useGetVehicleTargetConfigurationsQuery,
} from "@/apollo-client";
import VehicleForm from "./vehicle-form";

const defaultFilters: IVehicleTableFilters = {
  name: "",
};
export default function VehicleTable() {
  const [callDelete, callDeleteResponse] =
    useDeleteVehicleTargetConfigurationMutation();
  const { data, loading } = useGetVehicleTargetConfigurationsQuery({
    variables: {
      limit: 100,
      skip: 0,
    },
  });
  const table = useTable({
    defaultOrderBy: "createdAt",
    defaultOrder: "desc",
    defaultRowsPerPage: 10,
  });
  const [filters, setFilters] = useState(defaultFilters);
  const dataInTable = useMemo(
    () => data?.getVehicleTargetConfigurations ?? [],
    [data?.getVehicleTargetConfigurations]
  );
  const handleFilters = useCallback(
    (name: string, value: IVehicleTableFilterValue) => {
      table.onResetPage();
      setFilters((prevState) => ({
        ...prevState,
        [name]: value,
      }));
    },
    [table]
  );
  const handleDeleteRow = useCallback(
    (id: string) => {
      callDelete({
        variables: { id },
        refetchQueries: [GetVehicleTargetConfigurationsDocument],
      });
    },
    [callDelete]
  );
  const openVehicleForm = useBoolean();
  const [editVehicle, setEditVehicle] = useState(undefined);
  const handleEdit = useCallback(
    (item: any) => {
      setEditVehicle(item);
      openVehicleForm.onTrue();
    },
    [openVehicleForm]
  );
  return (
    <Card>
      <VehicleForm
        vehicle={editVehicle}
        opened={openVehicleForm.value}
        onClose={openVehicleForm.onFalse}
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
                  .map((row: any) => (
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
        disabled={loading}
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
