package com.service.post.service;

import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.stream.Collectors;

import org.springframework.grpc.server.service.GrpcService;

import com.service.post.BaseProfileResponse;
import com.service.post.BaseUserResponse;
import com.service.post.GetManyRequest;
import com.service.post.RestoreManyRequest;
import com.service.post.RestoreOneRequest;
import com.service.post.RestoredResponse;
import com.service.post.TopicAdminResponse;
import com.service.post.TopicsAdminResponse;
import com.service.post.UpdateTopicRequest;
import com.service.post.UpdatedResponse;
import com.service.post.common.SlugUtil;
import com.service.post.entity.TopicEntity;
import com.service.post.exceptions.AlreadyExistsException;
import com.service.post.exceptions.ResourceNotFoundException;
import com.service.post.grpc_clients.UserClient;
import com.service.post.CreateTopicRequest;
import com.service.post.CreatedResponse;
import com.service.post.DeleteManyRequest;
import com.service.post.DeleteOneRequest;
import com.service.post.DeletedResponse;
import com.service.post.PostServiceGrpc.PostServiceImplBase;
import com.service.post.repository.TopicRepository;
import com.service.user.UserPublicResponse;
import com.service.user.UsersPublicResponse;

import io.grpc.Status;
import io.grpc.StatusRuntimeException;
import io.grpc.stub.StreamObserver;
import jakarta.transaction.Transactional;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

@GrpcService
@RequiredArgsConstructor
@FieldDefaults(makeFinal = true, level = AccessLevel.PRIVATE)
public class PostServiceImpl extends PostServiceImplBase {
  TopicRepository topicRepository;
  UserClient userClient;

  @Override
  public void createTopic(CreateTopicRequest request, StreamObserver<CreatedResponse> responseObserver) {
    try {
      String slug = request.hasSlug() && !request.getSlug().isEmpty() ? request.getSlug()
          : SlugUtil.toSlug(request.getName());

      if (topicRepository.existsBySlug(slug)) {
        throw new AlreadyExistsException("slug đã tồn tại");
      }

      TopicEntity topic = TopicEntity.builder().name(request.getName()).slug(slug).createdById(request.getUserId())
          .updatedById(request.getUserId()).build();
      topicRepository.save(topic);

      CreatedResponse response = CreatedResponse.newBuilder().setId(topic.getId()).build();

      responseObserver.onNext(response);
      responseObserver.onCompleted();
    } catch (AlreadyExistsException e) {
      responseObserver.onError(Status.ALREADY_EXISTS.withDescription(e.getMessage()).asRuntimeException());
    } catch (Exception e) {
      responseObserver.onError(
          Status.INTERNAL.withDescription("tạo chủ đề bài viết thất bại: " + e.getMessage()).asRuntimeException());
    }
  }

  @Override
  public void getAllTopicsAdmin(GetManyRequest request, StreamObserver<TopicsAdminResponse> responseObserver) {
    try {
      List<TopicEntity> topics = topicRepository.findAll();

      if (topics.isEmpty()) {
        responseObserver.onNext(TopicsAdminResponse.newBuilder().build());
        responseObserver.onCompleted();
        return;
      }

      Set<String> userIdSet = new HashSet<>();
      for (TopicEntity t : topics) {
        userIdSet.add(t.getCreatedById());
        userIdSet.add(t.getUpdatedById());
      }
      List<String> userIds = new ArrayList<>(userIdSet);

      UsersPublicResponse usersRes = userClient.getUsersById(userIds);

      Map<String, UserPublicResponse> usersMap = usersRes.getUsersList().stream()
          .collect(Collectors.toMap(UserPublicResponse::getId, u -> u));

      TopicsAdminResponse.Builder responseBuilder = TopicsAdminResponse.newBuilder();
      for (TopicEntity t : topics) {
        TopicAdminResponse.Builder topicBuilder = TopicAdminResponse.newBuilder().setId(t.getId()).setName(t.getName())
            .setSlug(t.getSlug()).setCreatedAt(t.getCreatedAt().toString()).setUpdatedAt(t.getUpdatedAt().toString());

        if (usersMap.containsKey(t.getCreatedById())) {
          topicBuilder.setCreatedBy(toBaseUserResponse(usersMap.get(t.getCreatedById())));
        }

        if (usersMap.containsKey(t.getUpdatedById())) {
          topicBuilder.setUpdatedBy(toBaseUserResponse(usersMap.get(t.getUpdatedById())));
        }

        responseBuilder.addTopics(topicBuilder);
      }

      responseObserver.onNext(responseBuilder.build());
      responseObserver.onCompleted();
      return;
    } catch (StatusRuntimeException e) {
      responseObserver.onError(e);
    } catch (Exception e) {
      responseObserver.onError(
          Status.INTERNAL.withDescription("lấy danh sách chủ đề thất bại: " + e.getMessage()).asRuntimeException());
    }
  }

  @Override
  @Transactional
  public void updateTopic(UpdateTopicRequest request, StreamObserver<UpdatedResponse> responseObserver) {
    try {
      TopicEntity topic = topicRepository.findByIdAndDeletedTopicFalse(request.getId())
          .orElseThrow(() -> new ResourceNotFoundException("không tìm thấy chủ đề"));

      if (!topic.getName().equals(request.getName())) {
        topic.setName(request.getName());
      }
      if (!topic.getSlug().equals(request.getSlug())) {
        if (topicRepository.existsBySlug(request.getSlug())) {
          throw new AlreadyExistsException("slug đã tồn tại");
        }
        topic.setSlug(request.getSlug());
      }
      if (!topic.getUpdatedById().equals(request.getUserId())) {
        topic.setUpdatedById(request.getUserId());
      }

      topicRepository.save(topic);

      UpdatedResponse response = UpdatedResponse.newBuilder().setSuccess(true).build();

      responseObserver.onNext(response);
      responseObserver.onCompleted();
    } catch (AlreadyExistsException e) {
      responseObserver.onError(Status.ALREADY_EXISTS.withDescription(e.getMessage()).asRuntimeException());
    } catch (ResourceNotFoundException e) {
      responseObserver.onError(Status.NOT_FOUND.withDescription(e.getMessage()).asRuntimeException());
    } catch (Exception e) {
      responseObserver.onError(
          Status.INTERNAL.withDescription("cập nhật chủ đề bài viết thất bại: " + e.getMessage()).asRuntimeException());
    }
  }

  @Override
  public void deleteTopic(DeleteOneRequest request, StreamObserver<DeletedResponse> responseObserver) {
    try {
      TopicEntity topic = topicRepository.findByIdAndDeletedTopicFalse(request.getId())
          .orElseThrow(() -> new ResourceNotFoundException("không tìm thấy chủ đề bài viết"));

      topic.setDeletedTopic(true);

      if (!topic.getUpdatedById().equals(request.getUserId())) {
        topic.setUpdatedById(request.getUserId());
      }

      topicRepository.save(topic);

      DeletedResponse response = DeletedResponse.newBuilder().setSuccess(true).build();

      responseObserver.onNext(response);
      responseObserver.onCompleted();
    } catch (ResourceNotFoundException e) {
      responseObserver.onError(Status.NOT_FOUND.withDescription(e.getMessage()).asRuntimeException());
    } catch (Exception e) {
      responseObserver.onError(
          Status.INTERNAL.withDescription("chuyển chủ đề bài viết vào thùng rác thất bại: " + e.getMessage())
              .asRuntimeException());
    }
  }

  @Override
  @Transactional
  public void deleteTopics(DeleteManyRequest request, StreamObserver<DeletedResponse> responseObserver) {
    try {
      List<TopicEntity> topics = topicRepository.findAllByIdInAndDeletedTopicFalse(request.getIdsList());
      if (topics.size() != request.getIdsCount()) {
        throw new ResourceNotFoundException("Có chủ đề không tìm thấy");
      }

      topicRepository.updateIsDeletedAllById(request.getIdsList(), true, request.getUserId());

      DeletedResponse response = DeletedResponse.newBuilder().setSuccess(true).build();

      responseObserver.onNext(response);
      responseObserver.onCompleted();
    } catch (ResourceNotFoundException e) {
      responseObserver.onError(Status.NOT_FOUND.withDescription(e.getMessage()).asRuntimeException());
    } catch (Exception e) {
      responseObserver.onError(
          Status.INTERNAL.withDescription("chuyển danh sách chủ đề vào thùng rác thất bại: " + e.getMessage())
              .asRuntimeException());
    }
  }

  @Override
  public void restoreTopic(RestoreOneRequest request, StreamObserver<RestoredResponse> responseObserver) {
    try {
      TopicEntity topic = topicRepository.findByIdAndDeletedTopicTrue(request.getId())
          .orElseThrow(() -> new ResourceNotFoundException("không tìm thấy chủ đề bài viết"));

      topic.setDeletedTopic(false);

      if (!topic.getUpdatedById().equals(request.getUserId())) {
        topic.setUpdatedById(request.getUserId());
      }

      topicRepository.save(topic);

      RestoredResponse response = RestoredResponse.newBuilder().setSuccess(true).build();

      responseObserver.onNext(response);
      responseObserver.onCompleted();
    } catch (ResourceNotFoundException e) {
      responseObserver.onError(Status.NOT_FOUND.withDescription(e.getMessage()).asRuntimeException());
    } catch (Exception e) {
      responseObserver.onError(
          Status.INTERNAL.withDescription("khôi phục chủ đề bài viết thất bại: " + e.getMessage())
              .asRuntimeException());
    }
  }

  @Override
  @Transactional
  public void restoreTopics(RestoreManyRequest request, StreamObserver<RestoredResponse> responseObserver) {
    try {
      List<TopicEntity> topics = topicRepository.findAllByIdInAndDeletedTopicTrue(request.getIdsList());
      if (topics.size() != request.getIdsCount()) {
        throw new ResourceNotFoundException("Có chủ đề không tìm thấy");
      }

      topicRepository.updateIsDeletedAllById(request.getIdsList(), false, request.getUserId());

      RestoredResponse response = RestoredResponse.newBuilder().setSuccess(true).build();

      responseObserver.onNext(response);
      responseObserver.onCompleted();
    } catch (ResourceNotFoundException e) {
      responseObserver.onError(Status.NOT_FOUND.withDescription(e.getMessage()).asRuntimeException());
    } catch (Exception e) {
      responseObserver.onError(
          Status.INTERNAL.withDescription("khôi phục danh sách chủ đề thất bại: " + e.getMessage())
              .asRuntimeException());
    }
  }

  private BaseUserResponse toBaseUserResponse(UserPublicResponse u) {
    return BaseUserResponse.newBuilder().setId(u.getId()).setUsername(u.getUsername())
        .setProfile(BaseProfileResponse.newBuilder().setId(u.getProfile().getId())
            .setFirstName(u.getProfile().getFirstName()).setLastName(u.getProfile().getLastName()).build())
        .build();
  }
}
