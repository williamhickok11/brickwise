package com.brickwise.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.CorsRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;

/**
 * HTTP / CORS configuration.
 *
 * <p><b>C# analogy:</b> Like {@code builder.Services.AddCors(...)} and
 * {@code app.UseCors(...)} in ASP.NET Core {@code Program.cs}.
 */
@Configuration
public class WebConfig {

    @Bean
    public WebMvcConfigurer corsConfigurer() {
        return new WebMvcConfigurer() {
            @Override
            public void addCorsMappings(CorsRegistry registry) {
                registry.addMapping("/**")
                        .allowedOriginPatterns("*")
                        .allowedMethods("GET", "POST", "PUT", "DELETE", "OPTIONS")
                        .allowedHeaders("Accept", "Authorization", "Content-Type", "X-CSRF-Token")
                        .exposedHeaders("Link")
                        .allowCredentials(false)
                        .maxAge(300);
            }
        };
    }
}
