export type Operator = {
  name: string;
  id: string;
}

export type MainStore = {
  operators: Operator[];
  tripUpdates: TripUpdateEntry[];
  selectedOperator: string | null;
  isWaitingTripUpdates: boolean;
}

export enum MessageType {
  operators = 'operators',
  tripUpdates = 'tripUpdates'
}

export type TripUpdateParm = {
  operatorId: string;
}

export type OperatorsRequest = {
  requestType: MessageType;
}

export type TripUpdatesRequest = {
  requestType: MessageType;
  data: TripUpdateParm;
}

export type MessageRequest = OperatorsRequest | TripUpdatesRequest

export type OperatorsResponse = {
  responseType: MessageType;
  data: Operator[];
}

export type TripUpdatesResponse = {
  responseType: MessageType;
  data: any;
}

export type MessageResponse = OperatorsResponse | TripUpdatesResponse

export type TripUpdatesData = {
  operatorId: string;
  tripUpdates: TripUpdateEntry[];
}

export type TripUpdateEntry = {
  id: string;
  trip_update: TripUpdate;
}

export type TripUpdate = {
  trip: Trip;
  vehicle: Vehicle;
  stop_time_update: StopTimeEntry[];
  timestamp: number;
}


export type StopTime = { time: number; delay?: number; uncertainty?: number }
export type StopTimeEntry = {
  arrival?: StopTime;
  departure?: StopTime;
  schedule_relation_ship?: number;
  stop_sequence: number;
  stop_id: string;
}

export type Trip = {
  trip_id: string;
  route_id: string;
  direction_id: number;
  schedule_relation_ship?: number;
}

export type Vehicle = {
  id: string;
  label: string;
  license_plate: string;
}