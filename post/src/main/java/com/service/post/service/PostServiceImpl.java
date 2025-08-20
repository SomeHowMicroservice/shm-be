package com.service.post.service;

import org.springframework.grpc.server.service.GrpcService;

import com.service.post.protobuf.CreateTopicRequest;
import com.service.post.protobuf.CreatedResponse;
import com.service.post.protobuf.PostServiceGrpc.PostServiceImplBase;

import io.grpc.stub.StreamObserver;

@GrpcService
public class PostServiceImpl extends PostServiceImplBase {
  @Override
  public void createTopic(CreateTopicRequest request, StreamObserver<CreatedResponse> responseObserver) {
    System.out.println(request);
    CreatedResponse response = CreatedResponse.newBuilder()
        .setId("12345")
        .build();

    responseObserver.onNext(response);
    responseObserver.onCompleted();
  }
}
