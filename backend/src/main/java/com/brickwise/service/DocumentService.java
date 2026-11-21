package com.brickwise.service;

import com.brickwise.config.BrickwiseProperties;
import com.brickwise.error.ApiException;
import com.brickwise.model.Document;
import com.brickwise.repository.DocumentRepository;
import com.brickwise.repository.ResidenceRepository;
import org.springframework.core.io.FileSystemResource;
import org.springframework.core.io.Resource;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.io.InputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.StandardCopyOption;

/** One lease PDF per residence, stored on local disk (parity with Go DocumentService). */
@Service
public class DocumentService {
    private static final long MAX_BYTES = 10L << 20;

    private final DocumentRepository repo;
    private final ResidenceRepository residenceRepo;
    private final Path dir;

    public DocumentService(DocumentRepository repo, ResidenceRepository residenceRepo, BrickwiseProperties props) {
        this.repo = repo;
        this.residenceRepo = residenceRepo;
        this.dir = Path.of(props.getDocumentsDir());
    }

    public Document getByResidenceId(int residenceId) {
        residenceRepo.findById(residenceId).orElseThrow(ApiException::notFound);
        return repo.findByResidenceId(residenceId).orElseThrow(ApiException::notFound);
    }

    public Document upload(int residenceId, MultipartFile file) {
        residenceRepo.findById(residenceId).orElseThrow(ApiException::notFound);
        if (file == null || file.isEmpty()) {
            throw new ApiException("INVALID_FILE", "file is empty", 400);
        }
        if (file.getSize() > MAX_BYTES) {
            throw new ApiException("FILE_TOO_LARGE", "lease PDF must be 10MB or smaller", 400);
        }
        String name = file.getOriginalFilename() == null ? "lease.pdf" : file.getOriginalFilename();
        if (!name.toLowerCase().endsWith(".pdf")) {
            throw new ApiException("INVALID_FILE_TYPE", "only PDF files are allowed", 400);
        }
        String ct = file.getContentType() == null ? "application/pdf" : file.getContentType().toLowerCase();
        if (!ct.isBlank() && !ct.startsWith("application/pdf")) {
            throw new ApiException("INVALID_FILE_TYPE", "content type must be application/pdf", 400);
        }
        try {
            Files.createDirectories(dir);
            Path target = storagePath(residenceId);
            Path tmp = Path.of(target + ".tmp");
            try (InputStream in = file.getInputStream()) {
                Files.copy(in, tmp, StandardCopyOption.REPLACE_EXISTING);
            }
            Files.move(tmp, target, StandardCopyOption.REPLACE_EXISTING, StandardCopyOption.ATOMIC_MOVE);
        } catch (IOException e) {
            throw ApiException.internal("Failed to store lease PDF: " + e.getMessage());
        }
        Document doc = new Document();
        doc.setResidenceId(residenceId);
        doc.setFileName(name);
        doc.setContentType("application/pdf");
        doc.setSize(file.getSize());
        return repo.upsert(doc);
    }

    public Resource download(int residenceId) {
        getByResidenceId(residenceId);
        Path path = storagePath(residenceId);
        if (!Files.exists(path)) {
            throw ApiException.notFound();
        }
        return new FileSystemResource(path);
    }

    private Path storagePath(int residenceId) {
        return dir.resolve("residence-" + residenceId + ".pdf");
    }
}
