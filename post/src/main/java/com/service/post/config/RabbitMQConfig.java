package com.service.post.config;

import org.springframework.amqp.core.BindingBuilder;
import org.springframework.amqp.core.Binding;
import org.springframework.amqp.core.Queue;
import org.springframework.amqp.core.TopicExchange;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.amqp.support.converter.Jackson2JsonMessageConverter;
import org.springframework.amqp.rabbit.connection.ConnectionFactory;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class RabbitMQConfig {
  public static final String UPLOAD_QUEUE = "upload_queue";
  public static final String DELETE_QUEUE = "delete_queue";

  public static final String EXCHANGE = "image_exchange";

  public static final String UPLOAD_ROUTING_KEY = "image.upload";
  public static final String DELETE_ROUTING_KEY = "image.delete";

  @Bean
  public Queue uploadQueue() {
    return new Queue(UPLOAD_QUEUE);
  }

  @Bean
  public Queue deleteQueue() {
    return new Queue(DELETE_QUEUE);
  }

  @Bean
  public TopicExchange exchange() {
    return new TopicExchange(EXCHANGE);
  }

  @Bean
  public Binding bindingUploadQueue(Queue uploadQueue, TopicExchange exchange) {
    return BindingBuilder.bind(uploadQueue).to(exchange).with(UPLOAD_ROUTING_KEY);
  }

  @Bean
  public Binding bindingDeleteQueue(Queue deleteQueue, TopicExchange exchange) {
    return BindingBuilder.bind(deleteQueue).to(exchange).with(DELETE_ROUTING_KEY);
  }

  @Bean
  public Jackson2JsonMessageConverter jsonMessageConverter() {
    return new Jackson2JsonMessageConverter();
  }

  @Bean
  public RabbitTemplate rabbitTemplate(ConnectionFactory connectionFactory) {
    RabbitTemplate rabbitTemplate = new RabbitTemplate(connectionFactory);
    rabbitTemplate.setMessageConverter(jsonMessageConverter());
    return rabbitTemplate;
  }
}
