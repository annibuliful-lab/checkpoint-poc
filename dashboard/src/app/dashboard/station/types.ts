// ----------------------------------------------------------------------

import { StationLocation } from "@/apollo-client";

export type IStationTableFilterValue = string | Date | null;

export type IStationTableFilters = {
  name: string;
};

// ----------------------------------------------------------------------

export type IStationHistory = {
  StationTime: Date;
  paymentTime: Date;
  deliveryTime: Date;
  completionTime: Date;
  timeline: {
    title: string;
    time: Date;
  }[];
};

export type IStationShippingAddress = {
  fullAddress: string;
  phoneNumber: string;
};

export type IStationPayment = {
  cardType: string;
  cardNumber: string;
};

export type IStationDelivery = {
  shipBy: string;
  speedy: string;
  trackingNumber: string;
};

export type IStationCustomer = {
  id: string;
  name: string;
  email: string;
  avatarUrl: string;
  ipAddress: string;
};

export type IStationItem = StationLocation;
