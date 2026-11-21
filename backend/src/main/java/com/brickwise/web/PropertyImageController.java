package com.brickwise.web;

import com.brickwise.model.PropertyImage;
import com.brickwise.service.PropertyImageService;
import org.springframework.core.io.Resource;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import java.util.List;

@RestController
@RequestMapping("/api/v1/properties/{propertyId}/images")
public class PropertyImageController {
    private final PropertyImageService service;

    public PropertyImageController(PropertyImageService service) {
        this.service = service;
    }

    @GetMapping
    public List<PropertyImage> list(@PathVariable int propertyId) {
        return service.list(propertyId);
    }

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public PropertyImage upload(@PathVariable int propertyId, @RequestParam("file") MultipartFile file) {
        return service.upload(propertyId, file);
    }

    @GetMapping("/{imageId}/file")
    public ResponseEntity<Resource> file(@PathVariable int propertyId, @PathVariable int imageId) {
        PropertyImage meta = service.list(propertyId).stream()
                .filter(i -> i.getId() == imageId).findFirst()
                .orElseGet(() -> { throw com.brickwise.error.ApiException.notFound(); });
        Resource body = service.serve(propertyId, imageId);
        MediaType type = MediaType.parseMediaType(meta.getContentType());
        return ResponseEntity.ok().contentType(type).body(body);
    }

    @PutMapping("/{imageId}/main")
    public PropertyImage setMain(@PathVariable int propertyId, @PathVariable int imageId) {
        return service.setMain(propertyId, imageId);
    }

    @DeleteMapping("/{imageId}")
    public ResponseEntity<Void> delete(@PathVariable int propertyId, @PathVariable int imageId) {
        service.delete(propertyId, imageId);
        return ResponseEntity.noContent().build();
    }
}
