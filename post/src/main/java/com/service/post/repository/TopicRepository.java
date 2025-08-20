package com.service.post.repository;

import org.springframework.data.jpa.repository.JpaRepository;

import com.service.post.entity.TopicEntity;

public interface TopicRepository extends JpaRepository<TopicEntity, String> {
  
}
