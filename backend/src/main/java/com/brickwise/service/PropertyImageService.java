package com.brickwise.service;

import com.brickwise.config.BrickwiseProperties;
import com.brickwise.error.ApiException;
import com.brickwise.model.PropertyImage;
import com.brickwise.repository.PropertyImageRepository;
import com.brickwise.repository.PropertyRepository;
import org.springframework.core.io.FileSystemResource;
import org.springframework.core.io.Resource;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.io.InputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.StandardCopyOption;
import java.util.List;
import java.util.Locale;
import java.util.Set;

@Service
public class PropertyImageService {
    private static final long MAX_BYTES = 8L << 20;
    private static final Set<String> ALLOWED = Set.of(
            "image/jpeg", "image/png", "image/webp", "image/gif");

    private final PropertyImageRepository repo;
    private final PropertyRepository propertyRepo;
    private final Path dir;

    public PropertyImageService(PropertyImageRepository repo, PropertyRepository propertyRepo, BrickwiseProperties props) {
        this.repo = repo;
        this.propertyRepo = propertyRepo;
        this.dir = Path.of(props.getPropertyImagesDir());
    }

    public List<PropertyImage> list(int propertyId) {
        propertyRepo.findById(propertyId).orElseThrow(ApiException::notFound);
        return repo.listByPropertyId(propertyId);
    }

    public PropertyImage upload(int propertyId, MultipartFile file) {
        propertyRepo.findById(propertyId).orElseThrow(ApiException::notFound);
        if (file == null || file.isEmpty()) {
            throw new ApiException("INVALID_FILE", "file is empty", 400);
        }
        if (file.getSize() > MAX_BYTES) {
            throw new ApiException("FILE_TOO_LARGE", "image must be 8MB or smaller", 400);
        }
        String ct = file.getContentType() == null ? "" : file.getContentType().toLowerCase(Locale.ROOT);
        if (!ALLOWED.contains(ct)) {
            throw new ApiException("INVALID_FILE_TYPE", "image must be jpeg, png, webp, or gif", 400);
        }
        String name = file.getOriginalFilename() == null ? "image" : file.getOriginalFilename();
        boolean makeMain = repo.countByPropertyId(propertyId) == 0;
        PropertyImage meta = new PropertyImage();
        meta.setPropertyId(propertyId);
        meta.setFileName(name);
        meta.setContentType(ct);
        meta.setSize(file.getSize());
        meta.setMain(makeMain);
        PropertyImage saved = repo.create(meta);
        try {
            Files.createDirectories(dir);
            Path target = storagePath(saved.getId());
            try (InputStream in = file.getInputStream()) {
                Files.copy(in, target, StandardCopyOption.REPLACE_EXISTING);
            }
            repo.updateSize(saved.getId(), Files.size(target));
        } catch (IOException e) {
            repo.delete(saved.getId());
            throw ApiException.internal("Failed to store image: " + e.getMessage());
        }
        return repo.findById(saved.getId()).orElseThrow();
    }

    public Resource serve(int propertyId, int imageId) {
        PropertyImage img = repo.findById(imageId).orElseThrow(ApiException::notFound);
        if (img.getPropertyId() != propertyId) throw ApiException.notFound();
        Path path = storagePath(imageId);
        if (!Files.exists(path)) throw ApiException.notFound();
        return new FileSystemResource(path);
    }

    public PropertyImage setMain(int propertyId, int imageId) {
        PropertyImage img = repo.findById(imageId).orElseThrow(ApiException::notFound);
        if (img.getPropertyId() != propertyId) throw ApiException.notFound();
        repo.clearMain(propertyId);
        repo.setMain(imageId);
        return repo.findById(imageId).orElseThrow();
    }

    public void delete(int propertyId, int imageId) {
        PropertyImage img = repo.findById(imageId).orElseThrow(ApiException::notFound);
        if (img.getPropertyId() != propertyId) throw ApiException.notFound();
        repo.delete(imageId);
        try { Files.deleteIfExists(storagePath(imageId)); } catch (IOException ignored) {}
        if (img.isMain()) {
            List<PropertyImage> remaining = repo.listByPropertyId(propertyId);
            if (!remaining.isEmpty()) {
                repo.setMain(remaining.get(0).getId());
            }
        }
    }

    public void cleanupForProperty(int propertyId) {
        for (PropertyImage img : repo.listByPropertyId(propertyId)) {
            try { Files.deleteIfExists(storagePath(img.getId())); } catch (IOException ignored) {}
        }
    }

    private Path storagePath(int imageId) {
        return dir.resolve("image-" + imageId);
    }
}
