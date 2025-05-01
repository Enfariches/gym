import * as jspb from 'google-protobuf'

import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb'; // proto import: "google/protobuf/empty.proto"


export class GetAdminRequest extends jspb.Message {
  getAdminId(): string;
  setAdminId(value: string): GetAdminRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAdminRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetAdminRequest): GetAdminRequest.AsObject;
  static serializeBinaryToWriter(message: GetAdminRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAdminRequest;
  static deserializeBinaryFromReader(message: GetAdminRequest, reader: jspb.BinaryReader): GetAdminRequest;
}

export namespace GetAdminRequest {
  export type AsObject = {
    adminId: string,
  }
}

export class UpdateAdminRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAdminRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateAdminRequest): UpdateAdminRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateAdminRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAdminRequest;
  static deserializeBinaryFromReader(message: UpdateAdminRequest, reader: jspb.BinaryReader): UpdateAdminRequest;
}

export namespace UpdateAdminRequest {
  export type AsObject = {
  }
}

export class IsAdminRequest extends jspb.Message {
  getAdminId(): string;
  setAdminId(value: string): IsAdminRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IsAdminRequest.AsObject;
  static toObject(includeInstance: boolean, msg: IsAdminRequest): IsAdminRequest.AsObject;
  static serializeBinaryToWriter(message: IsAdminRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IsAdminRequest;
  static deserializeBinaryFromReader(message: IsAdminRequest, reader: jspb.BinaryReader): IsAdminRequest;
}

export namespace IsAdminRequest {
  export type AsObject = {
    adminId: string,
  }
}

export class DeleteAdminRequest extends jspb.Message {
  getAdminId(): string;
  setAdminId(value: string): DeleteAdminRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteAdminRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteAdminRequest): DeleteAdminRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteAdminRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteAdminRequest;
  static deserializeBinaryFromReader(message: DeleteAdminRequest, reader: jspb.BinaryReader): DeleteAdminRequest;
}

export namespace DeleteAdminRequest {
  export type AsObject = {
    adminId: string,
  }
}

export class IsAdminResponse extends jspb.Message {
  getIsadmin(): boolean;
  setIsadmin(value: boolean): IsAdminResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IsAdminResponse.AsObject;
  static toObject(includeInstance: boolean, msg: IsAdminResponse): IsAdminResponse.AsObject;
  static serializeBinaryToWriter(message: IsAdminResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IsAdminResponse;
  static deserializeBinaryFromReader(message: IsAdminResponse, reader: jspb.BinaryReader): IsAdminResponse;
}

export namespace IsAdminResponse {
  export type AsObject = {
    isadmin: boolean,
  }
}

export class Admin extends jspb.Message {
  getId(): string;
  setId(value: string): Admin;

  getName(): string;
  setName(value: string): Admin;

  getSurname(): string;
  setSurname(value: string): Admin;

  getEmail(): string;
  setEmail(value: string): Admin;

  getDepartament(): string;
  setDepartament(value: string): Admin;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Admin.AsObject;
  static toObject(includeInstance: boolean, msg: Admin): Admin.AsObject;
  static serializeBinaryToWriter(message: Admin, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Admin;
  static deserializeBinaryFromReader(message: Admin, reader: jspb.BinaryReader): Admin;
}

export namespace Admin {
  export type AsObject = {
    id: string,
    name: string,
    surname: string,
    email: string,
    departament: string,
  }
}

