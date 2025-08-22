package com.service.post.entity;

import jakarta.persistence.Column;
import jakarta.persistence.FetchType;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.Entity;
import jakarta.persistence.Table;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Entity
@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder
@Table(name = "images")
public class ImageEntity extends BaseEntity {
  @Column(nullable = false, length = 255)
  private String url;

  @Column(nullable = true, columnDefinition = "CHAR(24)")
  private String fileId;

  @Column(nullable = false)
  private int sortOrder;

  @ManyToOne(fetch = FetchType.LAZY)
  @JoinColumn(name = "post_id", nullable = false)
  private PostEntity post;
}
