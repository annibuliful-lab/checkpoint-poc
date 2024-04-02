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
import { IImeiTableFilterValue, IImeiTableFilters } from "./types";
import Scrollbar from "@/components/scrollbar";
import ImeiTableRow from "./table-row";
import { TABLE_HEAD } from "./const";
import { useBoolean } from "@/hooks/use-boolean";
import _ from "lodash";
import {
  GetMobileDeviceConfigurationsQuery,
  MobileDeviceConfiguration,
  useDeleteMobileDeviceConfigurationMutation,
} from "@/apollo-client";
import ImeiForm from "./imei-form";

const defaultFilters: IImeiTableFilters = {
  name: "",
};
type Props = {
  data?: GetMobileDeviceConfigurationsQuery;
  loading: boolean;
};
export default function ImeiTable({ data, loading }: Props) {
  const [deleteDevice, deleteDeviceResponse] =
    useDeleteMobileDeviceConfigurationMutation();
  const table = useTable({
    defaultOrderBy: "createdAt",
    defaultOrder: "desc",
    defaultRowsPerPage: 10,
  });
  const [filters, setFilters] = useState(defaultFilters);
  const dataInTable = useMemo(
    () => data?.getMobileDeviceConfigurations ?? [],
    [data]
  );
  const handleFilters = useCallback(
    (name: string, value: IImeiTableFilterValue) => {
      table.onResetPage();
      setFilters((prevState) => ({
        ...prevState,
        [name]: value,
      }));
    },
    [table]
  );
  const handleDeleteRow = useCallback(
    async (id: string) => {
      await deleteDevice({ variables: { id } });
    },
    [deleteDevice]
  );
  const openImeiForm = useBoolean();
  const [editImei, setEditImei] = useState<
    MobileDeviceConfiguration | undefined
  >(undefined);
  const handleEdit = useCallback(
    (item: MobileDeviceConfiguration) => {
      setEditImei(item);
      openImeiForm.onTrue();
    },
    [openImeiForm]
  );
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
                    <ImeiTableRow
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

      <ImeiForm
        opened={openImeiForm.value}
        onClose={openImeiForm.onFalse}
        dafaultValues={editImei}
      />

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
